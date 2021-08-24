package command_bus

import (
	"reflect"
	"runtime"
	"strings"
)

func getFrame(skipFrames int) runtime.Frame {
	// We need the frame at index skipFrames+2, since we never want runtime.Callers and getFrame
	targetFrameIndex := skipFrames + 2

	// Set size to targetFrameIndex+2 to ensure we have room for one more caller than we need
	programCounters := make([]uintptr, targetFrameIndex+2)
	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "unknown"}
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}

	return frame
}

// myCaller returns the caller of the function that called it :)
func myCaller() commandName {
	// Skip GetCallerFunctionName and the function to get the caller of
	return commandName(strings.TrimSuffix(getFrame(2).Function, ".func1"))
}

func getEventName(event Event) eventName {
	if event == nil {
		return ""
	}
	return eventName(reflect.TypeOf(event).String())
}

func getCommandName(command Command) commandName {
	return commandName(strings.TrimSuffix(runtime.FuncForPC(reflect.ValueOf(command).Pointer()).Name(), "-fm"))
}

func getFieldsList(event Event) []string {
	result := []string{}
	if event == nil {
		return result
	}
	indirect := reflect.Indirect(reflect.ValueOf(event))
	for i := 0; i < indirect.NumField(); i++ {
		result = append(result, indirect.Type().Field(i).Name)
	}
	return result
}
