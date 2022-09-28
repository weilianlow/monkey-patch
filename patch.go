package monkey_patch

import (
	"bou.ke/monkey"
	"context"
	"fmt"
	"github.com/weilianlow/monkey-patch/util"
	"reflect"
)

var (
	patchFuncMap map[string]patchFunc
)

type PatchMeta struct {
	Target      interface{}
	MethodName  string
	Replacement replacementFunc
}

type patchFunc struct {
	enabled  func(*Data)
	disabled func()
	isFunc   bool
}

type replacementFunc func(*Data) interface{}

func init() {
	patchFuncMap = make(map[string]patchFunc, 0)
}

func AddPatchFunc(patchName string, patchMeta PatchMeta) {
	logger := util.GetLogger()
	enabledFunc := func(data *Data) {
		var guard *monkey.PatchGuard
		guard = monkey.Patch(patchMeta.Target, patchMeta.Replacement(data))
		data.Guard = GuardFunc{
			Unpatch: guard.Unpatch,
			Restore: guard.Restore,
		}
	}
	disabledFunc := func() {
		monkey.Unpatch(patchMeta.Target)
	}
	if _, found := patchFuncMap[patchName]; found {
		logger.Warn(fmt.Sprintf("already added data{name: '%s'} as patch function", patchName))
	} else {
		logger.Debug(fmt.Sprintf("successfully added data{name: '%s'} as patch function", patchName))
		patchFuncMap[patchName] = patchFunc{
			enabled:  enabledFunc,
			disabled: disabledFunc,
			isFunc:   true,
		}
	}
}

func AddPatchInstanceMethod(patchName string, patchMeta PatchMeta) {
	logger := util.GetLogger()
	enabledFunc := func(data *Data) {
		var guard *monkey.PatchGuard
		guard = monkey.PatchInstanceMethod(reflect.TypeOf(patchMeta.Target), patchMeta.MethodName, patchMeta.Replacement(data))
		data.Guard = GuardFunc{
			Unpatch: guard.Unpatch,
			Restore: guard.Restore,
		}
	}
	disabledFunc := func() {
		monkey.UnpatchInstanceMethod(reflect.TypeOf(patchMeta.Target), patchMeta.MethodName)
	}
	if _, found := patchFuncMap[patchName]; found {
		logger.Warn(fmt.Sprintf("already added data{name: '%s'} as patch method", patchName))
	} else {
		logger.Debug(fmt.Sprintf("successfully added data{name: '%s'} as patch method", patchName))
		patchFuncMap[patchName] = patchFunc{
			enabled:  enabledFunc,
			disabled: disabledFunc,
		}
	}
}

func (d *DataList) MonkeyPatchByConfig(ctx context.Context) {
	d.monkeyPatch(ctx, d.Config)
}

func (d *DataList) MonkeyPatchByName(ctx context.Context, dataNames ...string) {
	d.monkeyPatch(ctx, dataNames)
}

func (d *DataList) monkeyPatch(ctx context.Context, dataNames []string) {
	logger := util.GetLoggerFromContext(ctx).WithField("func", "monkeyPatch")
	if len(dataNames) > 0 {
		for _, dataName := range dataNames {
			data, dataExist := d.DataMap[dataName]
			patchFunc, patchFuncExist := patchFuncMap[dataName]
			patchType := "function"
			if !patchFunc.isFunc {
				patchType = "method"
			}
			if dataExist && patchFuncExist {
				debugMsg := ""
				if data.Enabled {
					debugMsg = fmt.Sprintf("enabling patchType %s for data{name: '%s'}", patchType, dataName)
					patchFunc.enabled(&data)
				} else {
					debugMsg = fmt.Sprintf("disabling patchType %s for data{name: '%s'}", patchType, dataName)
					patchFunc.disabled()
				}
				if len(debugMsg) > 0 {
					logger.Debug(debugMsg)
				}
			} else {
				warningMsg := ""
				if !dataExist && !patchFuncExist {
					warningMsg = fmt.Sprintf("unable to retrieve data{name: '%s'} from data and %s map", dataName, patchType)
				} else if !dataExist {
					warningMsg = fmt.Sprintf("unable to retrieve data{name: '%s'} from data map", dataName)
				} else if !patchFuncExist {
					warningMsg = fmt.Sprintf("unable to retrieve data{name: '%s'} from %s map", dataName, patchType)
				}
				if len(warningMsg) > 0 {
					logger.Debug(warningMsg)
				}
			}
		}
	} else {
		for dataName, patchFunc := range patchFuncMap {
			if data, dataExist := d.DataMap[dataName]; dataExist {
				patchType := "function"
				if !patchFunc.isFunc {
					patchType = "method"
				}
				debugMsg := ""
				if data.Enabled {
					debugMsg = fmt.Sprintf("enabling patchType %s for data{name: '%s'}", patchType, dataName)
					patchFunc.enabled(&data)
				} else {
					debugMsg = fmt.Sprintf("disabling patchType %s for data{name: '%s'}", patchType, dataName)
					patchFunc.disabled()
				}
				if len(debugMsg) > 0 {
					logger.Debug(debugMsg)
				}
			} else {
				logger.Warn(fmt.Sprintf("unable to retrieve data{name: '%s'} from data map", dataName))
			}
		}
	}
}
