package generator

import (
	"fmt"
	"strings"
)

func (i *insn) genCodeRdVs2Vm() []string {
	combinations := i.combinations([]LMUL{1}, []SEW{8}, []bool{false, true})

	res := make([]string, 0, len(combinations))
	for _, c := range combinations {
		builder := strings.Builder{}
		builder.WriteString(i.gTestDataAddr())
		builder.WriteString(c.comment())

		vd := int(c.LMUL1)
		vs2 := int(c.LMUL1) * 2
		builder.WriteString(i.gWriteRandomData(LMUL(3)))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(0, c.LMUL1, SEW(8)))

		builder.WriteString(fmt.Sprintf("addi a0, a0, %d\n", 1*i.vlenb()))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vd, c.LMUL1, SEW(8)))

		builder.WriteString(fmt.Sprintf("addi a0, a0, %d\n", 1*i.vlenb()))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vs2, c.LMUL1, SEW(8)))

		builder.WriteString("# -------------- TEST BEGIN --------------\n")
		builder.WriteString(i.gVsetvli(c.Vl, c.SEW, c.LMUL))
		builder.WriteString(fmt.Sprintf("%s s0, v%d%s\n",
			i.Name, vs2, v0t(c.Mask)))
		builder.WriteString("# -------------- TEST END   --------------\n")

		builder.WriteString(i.gMoveScalarToVector("s0", vd, SEW(64)))
		builder.WriteString(i.gStoreRegisterGroupIntoData(vd, c.LMUL1, SEW(64)))
		builder.WriteString(i.gMagicInsn(vd))

		res = append(res, builder.String())
	}
	return res
}
