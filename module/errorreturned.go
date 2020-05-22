package module

import (
	"context"
	"path"

	"github.com/chaosblade-io/chaosblade-spec-go/channel"
	"github.com/chaosblade-io/chaosblade-spec-go/spec"
	"github.com/sirupsen/logrus"

	"github.com/chaosblade-io/chaosblade-exec-cplus/common"
)

type ErrorReturnedActionSpec struct {
	spec.BaseExpActionCommandSpec
}

func NewErrorReturnedActionSpec() spec.ExpActionCommandSpec {
	return &ErrorReturnedActionSpec{
		spec.BaseExpActionCommandSpec{
			ActionMatchers: []spec.ExpFlagSpec{},
			ActionFlags: []spec.ExpFlagSpec{
				&spec.ExpFlag{
					Name:     "returnValue",
					Desc:     "Value returned. If you want return null, set --returnValue null",
					Required: true,
				},
			},
			ActionExecutor: &ErrorReturnedExecutor{},
		},
	}
}

func (e ErrorReturnedActionSpec) Name() string {
	return "return"
}

func (e ErrorReturnedActionSpec) Aliases() []string {
	return []string{}
}

func (e ErrorReturnedActionSpec) ShortDesc() string {
	return "error returned"
}

func (e ErrorReturnedActionSpec) LongDesc() string {
	return "error returned"
}

type ErrorReturnedExecutor struct {
	channel spec.Channel
}

func (e *ErrorReturnedExecutor) Name() string {
	return "return"
}

func (e *ErrorReturnedExecutor) Exec(uid string, ctx context.Context, model *spec.ExpModel) *spec.Response {
	if _, ok := spec.IsDestroy(ctx); ok {
		return spec.ReturnFail(spec.Code[spec.ServerError], "illegal processing")
	}
	returnValue := model.ActionFlags["returnValue"]
	if returnValue == "" {
		return spec.ReturnFail(spec.Code[spec.IllegalParameters], "less necessary returnValue value")
	}
	// search pid by process name
	processName := model.ActionFlags["processName"]
	if processName == "" {
		return spec.ReturnFail(spec.Code[spec.IllegalParameters], "less necessary processName value")
	}
	processCtx := context.WithValue(context.Background(), channel.ExcludeProcessKey, "blade")
	localChannel := channel.NewLocalChannel()
	pids, err := localChannel.GetPidsByProcessName(processName, processCtx)
	if err != nil {
		logrus.Warnf("get pids by %s process name err, %v", processName, err)
	}
	if pids == nil || len(pids) == 0 {
		args := buildArgs([]string{
			model.ActionFlags["fileLocateAndName"],
			model.ActionFlags["forkMode"],
			model.ActionFlags["libLoad"],
			model.ActionFlags["breakLine"],
			returnValue,
			model.ActionFlags["initParams"],
		})
		return localChannel.Run(context.Background(), path.Join(common.GetScriptPath(), common.BreakAndReturnScript), args)
	} else {
		args := buildArgs([]string{
			pids[0],
			model.ActionFlags["forkMode"],
			"",
			"",
			model.ActionFlags["breakLine"],
			returnValue,
			model.ActionFlags["initParams"],
		})
		return localChannel.Run(context.Background(), path.Join(common.GetScriptPath(), common.BreakAndReturnAttachScript), args)
	}
}

func (e *ErrorReturnedExecutor) SetChannel(channel spec.Channel) {
	e.channel = channel
}