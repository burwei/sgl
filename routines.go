package sgl

import (
	"fmt"
	"io/ioutil"
	"log"
	"runtime"
	"strings"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func Init(windowWidth int, windowHeight int, windowTitle string) *glfw.Window {
	runtime.LockOSThread()
	window := InitGlfwAndOpenGL(windowWidth, windowHeight, windowTitle)
	return window
}

func Terminate() {
	glfw.Terminate()
}

func InitGlfwAndOpenGL(width int, height int, title string) *glfw.Window {
	// init GLFW
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	// init OpenGL
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	return window
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func MakeProgram(vertexShaderSource, fragmentShaderSource string) uint32 {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}
	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)
	gl.UseProgram(program)

	return program
}

func MakeProgramFromFile(vertPath string, fragPath string) uint32 {
	b, err := ioutil.ReadFile(vertPath)
	if err != nil {
		fmt.Print(err)
	}

	vShader := fmt.Sprintf("%s\x00", string(b))

	b, err = ioutil.ReadFile(fragPath)
	if err != nil {
		fmt.Print(err)
	}

	fShader := fmt.Sprintf("%s\x00", string(b))

	return MakeProgram(vShader, fShader)
}

func BeforeMainLoop(window *glfw.Window, vp *Viewpoint) {
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)
	keyCallback := glfw.KeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		// capture only Press and Repeat actions
		if action == glfw.Release {
			return
		}
		switch key {
		case glfw.KeyUp:
			vp.Eye = mgl32.Vec3{vp.Eye[0], vp.Eye[1] + 10, vp.Eye[2]}
			vp.Target = mgl32.Vec3{vp.Target[0], vp.Target[1] + 10, vp.Target[2]}
			vp.Camera = mgl32.LookAtV(vp.Eye, vp.Target, vp.Top)
		case glfw.KeyDown:
			vp.Eye = mgl32.Vec3{vp.Eye[0], vp.Eye[1] - 10, vp.Eye[2]}
			vp.Target = mgl32.Vec3{vp.Target[0], vp.Target[1] - 10, vp.Target[2]}
			vp.Camera = mgl32.LookAtV(vp.Eye, vp.Target, vp.Top)
		case glfw.KeyLeft:
			vp.Eye = mgl32.Vec3{vp.Eye[0] - 10, vp.Eye[1], vp.Eye[2]}
			vp.Target = mgl32.Vec3{vp.Target[0] - 10, vp.Target[1], vp.Target[2]}
			vp.Camera = mgl32.LookAtV(vp.Eye, vp.Target, vp.Top)
		case glfw.KeyRight:
			vp.Eye = mgl32.Vec3{vp.Eye[0] + 10, vp.Eye[1], vp.Eye[2]}
			vp.Target = mgl32.Vec3{vp.Target[0] + 10, vp.Target[1], vp.Target[2]}
			vp.Camera = mgl32.LookAtV(vp.Eye, vp.Target, vp.Top)
		case glfw.KeyO:
			vp.Eye = mgl32.Vec3{0, 0, 1000}
			vp.Target = mgl32.Vec3{0, 0, 0}
			vp.Camera = mgl32.LookAtV(vp.Eye, vp.Target, vp.Top)
		case glfw.KeyEscape:
			window.SetShouldClose(true)
		}
	})
	scrollCallback := glfw.ScrollCallback(func(w *glfw.Window, xpos, ypos float64) {
		forwardVec := vp.Target.Sub(vp.Eye).Mul(0.005)
		leftVec := vp.Top.Cross(forwardVec).Mul(2)
		forwardX := forwardVec[0] * float32(ypos)
		forwardY := forwardVec[1] * float32(ypos)
		forwardZ := forwardVec[2] * float32(ypos)
		vp.Eye = mgl32.Vec3{
			vp.Eye[0] + forwardX,
			vp.Eye[1] + forwardY,
			vp.Eye[2] + forwardZ,
		}
		vp.Target = mgl32.Vec3{
			vp.Target[0] + float32(xpos)*leftVec[0] + forwardX,
			vp.Target[1] + float32(xpos)*leftVec[1] + forwardY,
			vp.Target[2] + float32(xpos)*leftVec[2] + forwardZ,
		}
		vp.Camera = mgl32.LookAtV(vp.Eye, vp.Target, vp.Top)
	})
	window.SetKeyCallback(keyCallback)
	window.SetScrollCallback(scrollCallback)
}

func BeforeDrawing() {
	// Clear before redraw
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func AfterDrawing(window *glfw.Window) {
	// Maintenance
	window.SwapBuffers()
	glfw.PollEvents()
}
