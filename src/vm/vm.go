package vm

import (
	"ash/code"
	"ash/compiler"
	"ash/object"
	"fmt"
)

const StackSize = 2048

var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}

type VM struct {
	constants    []object.Object
	instructions code.Instructions
	stack        []object.Object
	sp           int
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,
		stack:        make([]object.Object, StackSize),
		sp:           0,
	}
}

func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip])

		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

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
			pos := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip = pos - 1

		case code.OpJumpNotTruthy:
			pos := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip += 2

			condition := vm.pop()
			if !isTruthy(condition) {
				ip = pos - 1
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

	rightType := right.Type()
	leftType := left.Type()

	if leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ {
		return vm.executeBinaryIntegerOp(op, left, right)
	}

	return fmt.Errorf("unsupported types for binary operation: %s %s",
		leftType, rightType)
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
	case False:
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
	default:
		return true
	}
}
