package vm

import (
	"ash/code"
	"ash/compiler"
	"ash/object"
	"fmt"
)

const StackSize = 2048
const GlobalsSize = 65536
const MaxFrames = 1024

var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}
var Null = &object.Null{}

type VM struct {
	constants   []object.Object
	stack       []object.Object
	sp          int
	globals     []object.Object
	frames      []*Frame
	framesIndex int
}

func New(bytecode *compiler.Bytecode) *VM {
	mainFn := &object.CompiledFunction{Instructions: bytecode.Instructions}
	mainFrame := NewFrame(mainFn, 0)

	frames := make([]*Frame, MaxFrames)
	frames[0] = mainFrame

	return &VM{
		constants: bytecode.Constants,

		stack: make([]object.Object, StackSize),
		sp:    0,

		globals: make([]object.Object, GlobalsSize),

		frames:      frames,
		framesIndex: 1,
	}
}

func NewWithGlobalsStore(bytecode *compiler.Bytecode, s []object.Object) *VM {
	vm := New(bytecode)
	vm.globals = s
	return vm
}

func (vm *VM) Run() error {
	var ins code.Instructions
	var op code.Opcode
	var ip int

	for vm.currFrame().ip < len(vm.currFrame().Instructions()) {
		vm.currFrame().ip++

		ip = vm.currFrame().ip
		ins = vm.currFrame().Instructions()
		op = code.Opcode(ins[ip])

		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(ins[ip+1:])
			vm.currFrame().ip += 2

			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}

		case code.OpPop:
			vm.pop()

		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
			err := vm.executeBinaryOp(op)
			if err != nil {
				return err
			}

		case code.OpTrue:
			err := vm.push(True)
			if err != nil {
				return err
			}

		case code.OpFalse:
			err := vm.push(False)
			if err != nil {
				return err
			}

		case code.OpNull:
			err := vm.push(Null)
			if err != nil {
				return err
			}

		case code.OpEqual, code.OpNotEqual, code.OpGreaterThan:
			err := vm.executeComparison(op)
			if err != nil {
				return err
			}

		case code.OpBang:
			err := vm.executeBangOp()
			if err != nil {
				return err
			}

		case code.OpMinus:
			err := vm.executeMinusOp()
			if err != nil {
				return err
			}

		case code.OpJump:
			pos := int(code.ReadUint16(ins[ip+1:]))
			vm.currFrame().ip = pos - 1

		case code.OpJumpNotTruthy:
			pos := int(code.ReadUint16(ins[ip+1:]))
			vm.currFrame().ip += 2

			condition := vm.pop()
			if !isTruthy(condition) {
				vm.currFrame().ip = pos - 1
			}

		case code.OpSetGlobal:
			idx := code.ReadUint16(ins[ip+1:])
			vm.currFrame().ip += 2

			vm.globals[idx] = vm.pop()

		case code.OpGetGlobal:
			idx := code.ReadUint16(ins[ip+1:])
			vm.currFrame().ip += 2

			err := vm.push(vm.globals[idx])
			if err != nil {
				return err
			}

		case code.OpSetLocal:
			localIdx := code.ReadUint16(ins[ip+1:])
			vm.currFrame().ip += 1

			frame := vm.currFrame()

			vm.stack[frame.basePointer+int(localIdx)] = vm.pop()

		case code.OpGetLocal:
			localIdx := code.ReadUint16(ins[ip+1:])
			vm.currFrame().ip += 1

			frame := vm.currFrame()

			err := vm.push(vm.stack[frame.basePointer+int(localIdx)])
			if err != nil {
				return err
			}

		case code.OpArray:
			n := int(code.ReadUint16(ins[ip+1:]))
			vm.currFrame().ip += 2

			arr := vm.buildArray(vm.sp-n, vm.sp)
			vm.sp -= n

			err := vm.push(arr)
			if err != nil {
				return err
			}

		case code.OpHash:
			n := int(code.ReadUint16(ins[ip+1:]))
			vm.currFrame().ip += 2

			hash, err := vm.buildHash(vm.sp-n, vm.sp)
			if err != nil {
				return err
			}
			vm.sp -= n

			err = vm.push(hash)
			if err != nil {
				return err
			}

		case code.OpIndex:
			i := vm.pop()
			l := vm.pop()

			err := vm.executeIndexExpression(l, i)
			if err != nil {
				return err
			}

		case code.OpCall:
			vm.currFrame().ip++
			fn, ok := vm.stack[vm.sp-1].(*object.CompiledFunction)
			if !ok {
				return fmt.Errorf("calling non-function")
			}
			frame := NewFrame(fn, vm.sp)
			vm.pushFrame(frame)
			vm.sp = frame.basePointer + fn.NumLocals

		case code.OpReturnValue:
			returnValue := vm.pop()

			frame := vm.popFrame()
			vm.sp = frame.basePointer - 1

			err := vm.push(returnValue)
			if err != nil {
				return err
			}

		case code.OpReturn:
			frame := vm.popFrame()
			vm.sp = frame.basePointer - 1

			err := vm.push(Null)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (vm *VM) LastPoppedStackElem() object.Object {
	return vm.stack[vm.sp]
}

func (vm *VM) push(o object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = o
	vm.sp++

	return nil
}

func (vm *VM) pop() object.Object {
	o := vm.stack[vm.sp-1]
	vm.sp--
	return o
}

func (vm *VM) executeBinaryOp(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	leftType := left.Type()
	rightType := right.Type()

	switch {
	case leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ:
		return vm.executeBinaryIntegerOp(op, left, right)
	case leftType == object.STRING_OBJ && rightType == object.STRING_OBJ:
		return vm.executeBinaryStringOp(op, left, right)
	default:
		return fmt.Errorf("unsupported types for binary operation: %s %s",
			leftType, rightType)
	}
}

func (vm *VM) executeBinaryIntegerOp(op code.Opcode, left, right object.Object) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	var result int64

	switch op {
	case code.OpAdd:
		result = leftValue + rightValue
	case code.OpSub:
		result = leftValue - rightValue
	case code.OpMul:
		result = leftValue * rightValue
	case code.OpDiv:
		result = leftValue / rightValue
	}

	return vm.push(&object.Integer{Value: result})
}

func (vm *VM) executeBinaryStringOp(op code.Opcode, left, right object.Object) error {
	if op != code.OpAdd {
		return fmt.Errorf("unknown string operator: %d", op)
	}

	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value

	return vm.push(&object.String{Value: leftValue + rightValue})
}

func (vm *VM) executeComparison(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		return vm.executeIntegerComparison(op, left, right)
	}

	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(left == right))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(left != right))
	default:
		return fmt.Errorf("unknown operator: %d (%s %s)",
			op, left.Type(), right.Type())
	}
}

func (vm *VM) executeIntegerComparison(op code.Opcode, left, right object.Object) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(leftValue == rightValue))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(leftValue != rightValue))
	case code.OpGreaterThan:
		return vm.push(nativeBoolToBooleanObject(leftValue > rightValue))
	default:
		return fmt.Errorf("unknown operator: %d", op)
	}
}

func (vm *VM) executeBangOp() error {
	op := vm.pop()

	switch op {
	case True:
		return vm.push(False)
	case False, Null:
		return vm.push(True)
	default:
		return vm.push(False)
	}
}

func (vm *VM) executeMinusOp() error {
	op := vm.pop()

	if op.Type() != object.INTEGER_OBJ {
		return fmt.Errorf("unsupported type for negation: %s", op.Type())
	}

	value := op.(*object.Integer).Value
	return vm.push(&object.Integer{Value: -value})
}

func (vm *VM) buildArray(start, end int) *object.Array {
	arr := make([]object.Object, end-start)

	for i := start; i < end; i++ {
		arr[i-start] = vm.stack[i]
	}

	return &object.Array{Elements: arr}
}

func (vm *VM) buildHash(start, end int) (*object.Hash, error) {
	pairs := make(map[object.HashKey]object.HashPair)

	for i := start; i < end; i += 2 {
		k := vm.stack[i]
		v := vm.stack[i+1]

		pair := object.HashPair{Key: k, Value: v}

		hashKey, ok := k.(object.Hashable)
		if !ok {
			return nil, fmt.Errorf("unusable as hash key: %s", k.Type())
		}

		pairs[hashKey.HashKey()] = pair
	}

	return &object.Hash{Pairs: pairs}, nil
}

func (vm *VM) executeIndexExpression(l, i object.Object) error {
	switch {
	case l.Type() == object.ARRRAY_OBJ && i.Type() == object.INTEGER_OBJ:
		return vm.executeArrayIndex(l, i)
	case l.Type() == object.HASH_OBJ:
		return vm.executeHashIndex(l, i)
	default:
		return fmt.Errorf("index operator not supported: %s", l.Type())
	}
}

func (vm *VM) executeArrayIndex(arr, idx object.Object) error {
	i := idx.(*object.Integer).Value
	obj := arr.(*object.Array)
	objLen := int64(len(obj.Elements) - 1)

	if i < 0 || i > objLen {
		return vm.push(Null)
	}

	return vm.push(obj.Elements[i])
}

func (vm *VM) executeHashIndex(hash, idx object.Object) error {
	obj := hash.(*object.Hash)

	key, ok := idx.(object.Hashable)
	if !ok {
		return fmt.Errorf("unusable as hash key: %s", idx.Type())
	}

	pair, ok := obj.Pairs[key.HashKey()]
	if !ok {
		return vm.push(Null)
	}

	return vm.push(pair.Value)
}

func (vm *VM) currFrame() *Frame {
	return vm.frames[vm.framesIndex-1]
}

func (vm *VM) pushFrame(f *Frame) {
	vm.frames[vm.framesIndex] = f
	vm.framesIndex++
}

func (vm *VM) popFrame() *Frame {
	vm.framesIndex--
	return vm.frames[vm.framesIndex]
}

func nativeBoolToBooleanObject(b bool) *object.Boolean {
	if b {
		return True
	}
	return False
}

func isTruthy(o object.Object) bool {
	switch o := o.(type) {
	case *object.Boolean:
		return o.Value

	case *object.Null:
		return false

	default:
		return true
	}
}
