package evaluator

import (
	"fmt"

	"github.com/vancanhuit/monkey/internal/ast"
	"github.com/vancanhuit/monkey/internal/object"
)

var (
	True  = &object.Boolean{Value: true}
	False = &object.Boolean{Value: false}
	Null  = &object.Null{}
)

func Eval(node ast.Node) object.Object {
	switch n := node.(type) {
	case *ast.Program:
		return evalProgram(n)
	case *ast.ExpressionStatement:
		return Eval(n.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: n.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(n.Value)
	case *ast.PrefixExpression:
		right := Eval(n.Right)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(n.Operator, right)
	case *ast.InfixExpression:
		left := Eval(n.Left)
		if isError(left) {
			return left
		}
		right := Eval(n.Right)
		if isError(right) {
			return right
		}
		return evalInfixExpression(n.Operator, left, right)
	case *ast.BlockStatement:
		return evalBlockStatement(n)
	case *ast.IfExpression:
		return evalIfExpression(n)
	case *ast.ReturnStatement:
		value := Eval(n.Value)
		if isError(value) {
			return value
		}
		return &object.ReturnValue{Value: value}
	}
	return nil
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ErrorObj
	}
	return false
}

func evalProgram(program *ast.Program) object.Object {
	var result object.Object

	for _, stmt := range program.Statements {
		result = Eval(stmt)

		switch obj := result.(type) {
		case *object.ReturnValue:
			return obj.Value
		case *object.Error:
			return obj
		}
	}

	return result
}

func evalBlockStatement(stmt *ast.BlockStatement) object.Object {
	var result object.Object

	for _, statement := range stmt.Statements {
		result = Eval(statement)

		if result != nil {
			t := result.Type()
			if t == object.ReturnValueObj || t == object.ErrorObj {
				return result
			}
		}
	}

	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return &object.Error{
			Message: fmt.Sprintf(
				"unknown operator: %s%s", operator, right.Type()),
		}
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case True:
		return False
	case False:
		return True
	case Null:
		return True
	default:
		return False
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.IntegerObj {
		return &object.Error{
			Message: fmt.Sprintf("unknown operator: -%s", right.Type()),
		}
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(
	operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.IntegerObj && right.Type() == object.IntegerObj:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return &object.Error{
			Message: fmt.Sprintf("type mismatch: %s %s %s", left.Type(), operator, right.Type()),
		}
	default:
		return &object.Error{
			Message: fmt.Sprintf("unknown operator: %s %s %s", left.Type(), operator, right.Type()),
		}
	}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return True
	}
	return False
}

func evalIntegerInfixExpression(
	operator string, left, right object.Object) object.Object {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftValue + rightValue}
	case "-":
		return &object.Integer{Value: leftValue - rightValue}
	case "*":
		return &object.Integer{Value: leftValue * rightValue}
	case "/":
		return &object.Integer{Value: leftValue / rightValue}
	case "<":
		return nativeBoolToBooleanObject(leftValue < rightValue)
	case ">":
		return nativeBoolToBooleanObject(leftValue > rightValue)
	case "==":
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case "!=":
		return nativeBoolToBooleanObject(leftValue != rightValue)
	default:
		return &object.Error{
			Message: fmt.Sprintf("unknown operator: %s %s %s", left.Type(), operator, right.Type()),
		}
	}
}

func evalIfExpression(expr *ast.IfExpression) object.Object {
	condition := Eval(expr.Condition)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(expr.Consequence)
	} else if expr.Alternative != nil {
		return Eval(expr.Alternative)
	}

	return Null
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case Null:
		return false
	case True:
		return true
	case False:
		return false
	default:
		return true
	}
}
