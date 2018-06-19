package trcapi_test

import (
	"context"
	"encoding/json"
	"io"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/rvolosatovs/turtlitto/pkg/api"
	"github.com/rvolosatovs/turtlitto/pkg/api/apitest"
	. "github.com/rvolosatovs/turtlitto/pkg/trcapi"
	"github.com/rvolosatovs/turtlitto/pkg/trcapi/trctest"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

//Test_items: Connect(), SendHandshake(), SendState(), State(), SubscribeStateChanges() in conn.go
//Input_spec: -
//Output_spec: Pass or fail
//Envir_needs: -
func TestState(t *testing.T) {
	for i, tc := range []struct {
		Expected *api.State
	}{
		{
			Expected: &api.State{
				Turtles: map[string]*api.TurtleState{
					"1": {
						BatteryVoltage: apitest.Uint8Ptr(42),
					},
					"2": {
						HomeGoal: api.HomeGoalBlue,
					},
					"3": {},
					"4": {},
					"5": {},
					"6": {},
				},
			},
		},
		{
			Expected: &api.State{
				Command: api.CommandBallHandlingDemo,
				Turtles: map[string]*api.TurtleState{
					"1": {
						HomeGoal: api.HomeGoalBlue,
					},
					"2": {},
					"3": {},
					"4": {},
					"5": {},
					"6": {},
				},
			},
		},
		{
			Expected: &api.State{
				Command: api.CommandCornerMagenta,
				Turtles: map[string]*api.TurtleState{
					"1": {},
					"2": {},
					"3": {},
					"4": {},
					"5": {},
					"6": {},
				},
			},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			a := assert.New(t)

			srrsIn, trcOut := io.Pipe()
			trcIn, srrsOut := io.Pipe()

			trc := trctest.Connect(trcOut, trcIn,
				trctest.WithHandler(api.MessageTypeHandshake, trctest.DefaultHandshakeHandler),
			)

			wg := &sync.WaitGroup{}
			wg.Add(3)

			go func() {
				defer wg.Done()

				for err := range trc.Errors() {
					panic(errors.Wrap(err, "TRC error"))
				}
			}()

			go func() {
				defer wg.Done()

				err := trc.SendHandshake(&api.Handshake{Version: DefaultVersion})
				a.Nil(err)
			}()

			conn, err := Connect(DefaultVersion, srrsOut, srrsIn)
			a.Nil(err)

			go func() {
				defer wg.Done()

				for err := range conn.Errors() {
					panic(errors.Wrap(err, "SRRS error"))
				}
			}()

			ctx := context.Background()

			st := conn.State(ctx)
			a.Equal(&api.State{
				Turtles: map[string]*api.TurtleState{
					"1": {},
					"2": {},
					"3": {},
					"4": {},
					"5": {},
					"6": {},
				},
			}, st)

			ch, closeFn, err := conn.SubscribeStateChanges(ctx)
			a.NoError(err)
			a.NotNil(closeFn)

			err = trc.SendState(tc.Expected)
			a.NoError(err)

			select {
			case <-ch:
			case <-time.After(time.Second):
				t.Error("No update received")
				t.FailNow()
			}

			st = conn.State(ctx)
			a.Equal(tc.Expected, st)

			a.NotPanics(func() { closeFn() })

			select {
			case <-ch:
			case <-time.After(time.Second):
				t.Error("Subscription channel not closed")
				t.FailNow()
			}

			err = conn.Close()
			a.NoError(err)

			err = trc.Close()
			a.NoError(err)

			err = trcIn.Close()
			a.NoError(err)

			err = srrsIn.Close()
			a.NoError(err)

			wg.Wait()
		})
	}
}

//Test_items: Connect(), SendHandshake(), SendState(), State(), SubscribeStateChanges() in conn.go
//Input_spec: -
//Output_spec: Pass or fail
//Envir_needs: -
func TestSetState(t *testing.T) {
	logger := zap.L()

	for i, tc := range []struct {
		Input  *api.State
		Output *api.State
	}{
		{
			Input: &api.State{
				Turtles: map[string]*api.TurtleState{
					"1": {
						BatteryVoltage: apitest.Uint8Ptr(42),
					},
					"2": {
						HomeGoal: api.HomeGoalBlue,
					},
					"3": {},
					"4": {},
					"5": {},
					"6": {},
				},
			},
			Output: &api.State{
				Turtles: map[string]*api.TurtleState{
					"1": {},
					"2": {
						HomeGoal: api.HomeGoalBlue,
					},
					"3": {},
					"4": {},
					"5": {},
					"6": {},
				},
			},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			logger := logger.With(zap.String("test", t.Name()))

			a := assert.New(t)

			srrsIn, trcOut := io.Pipe()
			trcIn, srrsOut := io.Pipe()

			trc := trctest.Connect(trcOut, trcIn,
				trctest.WithHandler(api.MessageTypeHandshake, trctest.DefaultHandshakeHandler),
				trctest.WithHandler(api.MessageTypeState, func(msg *api.Message) (*api.Message, error) {
					a.Nil(msg.ParentID)
					a.NotEmpty(msg.MessageID)

					in := &api.State{}
					err := json.Unmarshal(msg.Payload, in)
					a.NoError(err)
					a.Equal(tc.Input, in)

					b, err := json.Marshal(tc.Output)
					if err != nil {
						panic(err)
					}
					return api.NewMessage(api.MessageTypeState, b, &msg.MessageID), nil
				}),
			)

			wg := &sync.WaitGroup{}
			wg.Add(3)

			go func() {
				defer wg.Done()

				for err := range trc.Errors() {
					panic(errors.Wrap(err, "TRC error"))
				}
			}()

			go func() {
				defer wg.Done()

				err := trc.SendHandshake(&api.Handshake{Version: DefaultVersion})
				a.Nil(err)
			}()

			conn, err := Connect(DefaultVersion, srrsOut, srrsIn)
			a.Nil(err)

			go func() {
				defer wg.Done()

				for err := range conn.Errors() {
					panic(errors.Wrap(err, "SRRS error"))
				}
			}()

			ctx := context.Background()

			st := conn.State(ctx)
			a.Equal(&api.State{
				Turtles: map[string]*api.TurtleState{
					"1": {},
					"2": {},
					"3": {},
					"4": {},
					"5": {},
					"6": {},
				},
			}, st)

			ch, closeFn, err := conn.SubscribeStateChanges(ctx)
			a.NoError(err)
			a.NotNil(closeFn)

			logger.Debug("Setting state...")
			err = conn.SetState(ctx, tc.Input)
			a.NoError(err)

			logger.Debug("Waiting for state update...")
			select {
			case <-ch:
			case <-time.After(time.Second):
				t.Error("No update received")
				t.FailNow()
			}

			logger.Debug("Querying the state...")
			st = conn.State(ctx)
			a.Equal(tc.Output, st)

			a.NotPanics(func() { closeFn() })

			select {
			case <-ch:
			case <-time.After(time.Second):
				t.Error("Subscription channel not closed")
				t.FailNow()
			}

			err = conn.Close()
			a.NoError(err)

			err = trc.Close()
			a.NoError(err)

			err = trcIn.Close()
			a.NoError(err)

			err = srrsIn.Close()
			a.NoError(err)

			wg.Wait()
		})
	}
}
