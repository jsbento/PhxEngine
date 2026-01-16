package init

import (
	"log"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
)

func InitOpenGL(debug bool) error {
	if err := gl.Init(); err != nil {
		return err
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	renderer := gl.GoStr(gl.GetString(gl.RENDERER))
	vendor := gl.GoStr(gl.GetString(gl.VENDOR))
	log.Printf("OpenGL version: %s", version)
	log.Printf("OpenGL renderer: %s", renderer)
	log.Printf("OpenGL vendor: %s", vendor)

	if debug {
		gl.Enable(gl.DEBUG_OUTPUT)
		gl.Enable(gl.DEBUG_OUTPUT_SYNCHRONOUS)
		gl.DebugMessageCallback(
			func(source, gltype, id uint32, severity uint32, length int32, message string, userParam unsafe.Pointer) {
				log.Printf(
					"OpenGL debug: severity=%s source=%s type=%s id=%d message=%s",
					debugSeverityString(severity),
					debugSourceString(source),
					debugTypeString(gltype),
					id,
					message,
				)
				if severity == gl.DEBUG_SEVERITY_HIGH {
					panic("OpenGL high severity error")
				}
			},
			nil,
		)
		gl.DebugMessageControl(gl.DONT_CARE, gl.DONT_CARE, gl.DONT_CARE, 0, nil, true)
	}

	return nil
}

func debugSourceString(source uint32) string {
	switch source {
	case gl.DEBUG_SOURCE_API:
		return "api"
	case gl.DEBUG_SOURCE_WINDOW_SYSTEM:
		return "window_system"
	case gl.DEBUG_SOURCE_SHADER_COMPILER:
		return "shader_compiler"
	case gl.DEBUG_SOURCE_THIRD_PARTY:
		return "third_party"
	case gl.DEBUG_SOURCE_APPLICATION:
		return "application"
	case gl.DEBUG_SOURCE_OTHER:
		return "other"
	default:
		return "unknown"
	}
}

func debugTypeString(gltype uint32) string {
	switch gltype {
	case gl.DEBUG_TYPE_ERROR:
		return "error"
	case gl.DEBUG_TYPE_DEPRECATED_BEHAVIOR:
		return "deprecated_behavior"
	case gl.DEBUG_TYPE_UNDEFINED_BEHAVIOR:
		return "undefined_behavior"
	case gl.DEBUG_TYPE_PORTABILITY:
		return "portability"
	case gl.DEBUG_TYPE_PERFORMANCE:
		return "performance"
	case gl.DEBUG_TYPE_MARKER:
		return "marker"
	case gl.DEBUG_TYPE_PUSH_GROUP:
		return "push_group"
	case gl.DEBUG_TYPE_POP_GROUP:
		return "pop_group"
	case gl.DEBUG_TYPE_OTHER:
		return "other"
	default:
		return "unknown"
	}
}

func debugSeverityString(severity uint32) string {
	switch severity {
	case gl.DEBUG_SEVERITY_HIGH:
		return "high"
	case gl.DEBUG_SEVERITY_MEDIUM:
		return "medium"
	case gl.DEBUG_SEVERITY_LOW:
		return "low"
	case gl.DEBUG_SEVERITY_NOTIFICATION:
		return "notification"
	default:
		return "unknown"
	}
}
