package example

import(
	"mylib/src/public"
	"mylib/src/module/bignum_manager"
)

func Example_bitnum(){

	test_1 := bignum_manager.Calc("( ", 1, "+", 2.5, " ) + 0x11*( 0x10 % ", 5, ")")
	test_2 := bignum_manager.Calc("( 0xfffffffffffff * 0xeeeeeeeeeeeeee) / 0xaaaaaaa + 1234567890987654321")
	test_3 := bignum_manager.Calc_Keep_Point("15924006418246166 / 17180078638064511399535 / 1000000000000000000 * 4000")

	public.DBG_LOG("calc result: ", test_1)
	public.DBG_LOG("calc result: ", test_2)
	public.DBG_LOG("calc result: ", test_3)
}

