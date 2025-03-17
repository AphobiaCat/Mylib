package bignum_manager

import (
	"mylib/src/public"
	"math/big"

	"fmt"
	"strconv"
)

type calc_item struct{
	num				*big.Float
	op_type			byte
	op_type_level	int
	is_op			bool
}

func (ci calc_item) is_operator()bool{
	return ci.is_op
}

type mid_calc struct{
	item 	[]calc_item
}

func (mc *mid_calc) add_op(op_type byte){

	op_type_level := 0

	if op_type == '+' || op_type == '-'{
		op_type_level = 2	
	}else if op_type == '/' || op_type == '*' || op_type == '%'{
		op_type_level = 1
	}else if op_type == ')'{
		op_type_level = 3
	}else{	//'('
		op_type_level = 4
	}

	mc.item = append(mc.item, calc_item{op_type: op_type, is_op: true, op_type_level: op_type_level})
}

func (mc *mid_calc) add_num(num *big.Float){
	mc.item = append(mc.item, calc_item{num: num, is_op: false})
}

func (mc *mid_calc) add_other(nums mid_calc){
	mc.item = append(mc.item, nums.item...)
}

func (mc *mid_calc) add_other_revers(nums []calc_item){

	for i:= len(nums) - 1; i >= 0; i--{	
		mc.item = append(mc.item, nums[i])
	}
}


func (mc *mid_calc) change_2_front(){
	var result mid_calc

	var op_stack []calc_item
	op_stack_index := -1

	for i := len(mc.item) - 1; i >= 0; i--{

		if mc.item[i].is_operator(){
		
			if op_stack_index == -1 || mc.item[i].op_type_level == 3{				

				op_stack = append(op_stack, mc.item[i])
				op_stack_index++
			}else if mc.item[i].op_type_level == 4{

				var k int
				for k = op_stack_index; k >= 0; k--{
					if op_stack[k].op_type_level == 3{
						break
					}
				}

				if k < 0{
					panic("k can't low than zero")
				}

				result.add_other_revers(op_stack[k + 1:])
				
				op_stack = op_stack[:k]
				op_stack_index = k - 1
			}else{				
			
				var k int
				for k = op_stack_index; k >= 0; k--{
					if op_stack[k].op_type_level < mc.item[i].op_type_level{
						result.add_op(op_stack[k].op_type)
					}else{
						break
					}
				}

				if k + 1 < 0{
					panic("k + 1 can't low than zero")
				}

				op_stack = op_stack[:k + 1]
				op_stack = append(op_stack, mc.item[i])

				op_stack_index = k + 1
			}
		}else{
			result.add_num(mc.item[i].num)
		}
	}

	result.add_other_revers(op_stack)

	mc.item = result.item
}

func (mc mid_calc) calc() string{
	var nums_array []*big.Float
	var array_index int = -1

	for _, val := range mc.item{
	
		if val.is_operator() && array_index - 1 >= 0{
			switch val.op_type{
			case '+':
				sum := new(big.Float).Add(nums_array[array_index], nums_array[array_index - 1])
				nums_array[array_index - 1] = sum
				nums_array = nums_array[:array_index]
				array_index -= 1
							
			case '-':
				sum := new(big.Float).Sub(nums_array[array_index], nums_array[array_index - 1])
				nums_array[array_index - 1] = sum
				nums_array = nums_array[:array_index]
				array_index -= 1
			
			case '*':
				sum := new(big.Float).Mul(nums_array[array_index], nums_array[array_index - 1])
				nums_array[array_index - 1] = sum
				nums_array = nums_array[:array_index]
				array_index -= 1
			
			case '/':
				sum := new(big.Float).Quo(nums_array[array_index], nums_array[array_index - 1])
				nums_array[array_index - 1] = sum
				nums_array = nums_array[:array_index]
				array_index -= 1
			
			case '%':
				dividend	:= nums_array[array_index]
				divisor		:= nums_array[array_index - 1]

				dividendInt, _ := dividend.Int(nil)
				divisorInt, _ := divisor.Int(nil)

				remainder := new(big.Int).Mod(dividendInt, divisorInt)
				remainderFloat := new(big.Float).SetInt(remainder)

				nums_array[array_index - 1] = remainderFloat
				nums_array = nums_array[:array_index]
				array_index -= 1
			}
		}else{
			nums_array	= append(nums_array, val.num)
			array_index	+= 1
		}
	}

	return nums_array[0].String()
}

func (mc mid_calc) print(){

	var print_str string

	for _, val := range mc.item{
		if val.is_operator(){
			print_str += fmt.Sprintf("%c ", val.op_type)	
		}else{
			print_str += fmt.Sprintf("%s ", val.num.String())
		}
	}

	public.DBG_LOG("calc:", print_str)
}

func process_string_calc(calc_str string) mid_calc {

	var ret mid_calc

	len_of_calc := len(calc_str)

	for i := 0; i < len_of_calc; i++ {
		val := calc_str[i]

		if val == '(' || val == ')' || val == '+' || val == '-' || val == '*' || val == '/' || val == '%' {
			ret.add_op(val)
		} else if val == ' ' {
			continue
		} else if (val >= '0' && val <= '9') {

			var k int
			for k = i; k < len_of_calc; k++ {
				tmp_val := calc_str[k]
				// expense Symbol
				if (tmp_val >= '0' && tmp_val <= '9') ||
					(tmp_val >= 'a' && tmp_val <= 'f') ||
					(tmp_val >= 'A' && tmp_val <= 'F') ||
					tmp_val == 'x' ||
					tmp_val == '.' ||
					tmp_val == 'e' ||  // since calc, even if include by hex
					tmp_val == 'E' ||  // since calc, even if include by hex
					tmp_val == '+' ||  // allow index symbol
					tmp_val == '-' {   // allow index symbol
					continue
				} else {
					break
				}
			}

			num_str := calc_str[i:k]
			i = k - 1

			tmp_big_num := new(big.Float)
			_, _, err := tmp_big_num.Parse(num_str, 0)
			if err != nil {
				public.DBG_ERR("parse num error:", err)
				continue
			}

			ret.add_num(tmp_big_num)
		}
	}

	return ret
}

func make_calc_mid(calc_item ...interface{}) mid_calc {

	var ret mid_calc

	for _, val := range calc_item{
		switch t := val.(type) {
			case string:
				calc_str_tmp := val.(string)
				num_tmp := process_string_calc(calc_str_tmp)
				ret.add_other(num_tmp)

			case float32:
				tmp_num := float64(val.(float32))
				ret.add_num(big.NewFloat(tmp_num))

			case float64:
				tmp_num := val.(float64)
				ret.add_num(big.NewFloat(tmp_num))

			case int:
				tmp_num := float64(val.(int))		
				ret.add_num(big.NewFloat(tmp_num))

			case uint:
				tmp_num := float64(val.(uint))		
				ret.add_num(big.NewFloat(tmp_num))

			case uint32:
				tmp_num := float64(val.(uint32))		
				ret.add_num(big.NewFloat(tmp_num))

			case int64:
				tmp_num := float64(val.(int64))		
				ret.add_num(big.NewFloat(tmp_num))

			case uint64:
				tmp_num := float64(val.(uint64))		
				ret.add_num(big.NewFloat(tmp_num))

			default:
				public.DBG_ERR("No support param type:", t)
        }
	}

	return ret
}

func Byte_2_Num_Str(byte_array [32]byte)string{
	hashInt := new(big.Int).SetBytes(byte_array[:])
	return hashInt.String()	
}

func Calc(calc_item ...interface{}) string{

	defer func(){
		if err := recover(); err != nil{
			public.DBG_ERR("err:", err)
		}
	}()

	num_mid := make_calc_mid(calc_item...)

	//num_mid.print()

	num_mid.change_2_front()

	//num_mid.print()

	return num_mid.calc()
}

func Calc_Keep_Point(calc_item ...interface{}) string{
	ret := Calc(calc_item...)

	floatVal, err := strconv.ParseFloat(ret, 64)
	if err != nil {
		public.DBG_LOG("ParseFloat err:", err)
		return ""
	}
	
	return strconv.FormatFloat(floatVal, 'f', -1, 64)
}

