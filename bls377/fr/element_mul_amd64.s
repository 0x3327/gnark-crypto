// Code generated by goff (v0.2.1) DO NOT EDIT

#include "textflag.h"

// func MulAssignElement(res,y *Element)
// montgomery multiplication of res by y
// stores the result in res
TEXT ·MulAssignElement(SB), NOSPLIT, $0-16
    
    // dereference our parameters
    MOVQ res+0(FP), DI
    MOVQ y+8(FP), R8
    
    // check if we support adx and mulx
    CMPB ·supportAdx(SB), $1
    JNE no_adx
    
    // the algorithm is described here
    // https://hackmd.io/@zkteam/modular_multiplication
    // however, to benefit from the ADCX and ADOX carry chains
    // we split the inner loops in 2:
    // for i=0 to N-1
    // 		for j=0 to N-1
    // 		    (A,t[j])  := t[j] + a[j]*b[i] + A
    // 		m := t[0]*q'[0] mod W
    // 		C,_ := t[0] + m*q[0]
    // 		for j=1 to N-1
    // 		    (C,t[j-1]) := t[j] + m*q[j] + C
    // 		t[N-1] = C + A
    
    // ---------------------------------------------------------------------------------------------
    // outter loop 0
    
    // clear up the carry flags
    XORQ R9 , R9
    
    // DX = y[0]
    MOVQ 0(R8), DX
    
    // for j=0 to N-1
    //    (A,t[j])  := t[j] + x[j]*y[i] + A
    
    MULXQ 0(DI), CX, BX
    
    MULXQ 8(DI), AX, BP
    ADOXQ AX, BX
    
    MULXQ 16(DI), AX, SI
    ADOXQ AX, BP
    
    MULXQ 24(DI), AX, R9
    ADOXQ AX, SI
    
    // add the last carries to R9
    MOVQ $0, DX
    ADCXQ DX, R9
    ADOXQ DX, R9
    
    // m := t[0]*q'[0] mod W
    MOVQ $0x0a117fffffffffff, DX
    MULXQ CX,R10, DX
    
    // clear the carry flags
    XORQ DX, DX
    
    // C,_ := t[0] + m*q[0]
    MOVQ $0x0a11800000000001, DX
    MULXQ R10, AX, DX
    ADCXQ CX ,AX
    MOVQ DX, CX
    
    // for j=1 to N-1
    //    (C,t[j-1]) := t[j] + m*q[j] + C
    MOVQ $0x59aa76fed0000001, DX
    ADCXQ  BX, CX
    MULXQ R10, AX, BX
    ADOXQ AX, CX
    MOVQ $0x60b44d1e5c37b001, DX
    ADCXQ  BP, BX
    MULXQ R10, AX, BP
    ADOXQ AX, BX
    MOVQ $0x12ab655e9a2ca556, DX
    ADCXQ  SI, BP
    MULXQ R10, AX, SI
    ADOXQ AX, BP
    MOVQ $0, AX
    ADCXQ AX, SI
    ADOXQ R9, SI
    
    // ---------------------------------------------------------------------------------------------
    // outter loop 1
    
    // clear up the carry flags
    XORQ R9 , R9
    
    // DX = y[1]
    MOVQ 8(R8), DX
    
    // for j=0 to N-1
    //    (A,t[j])  := t[j] + x[j]*y[i] + A
    
    MULXQ 0(DI), AX, R9
    ADOXQ AX, CX
    
    ADCXQ R9, BX
    MULXQ 8(DI), AX, R9
    ADOXQ AX, BX
    
    ADCXQ R9, BP
    MULXQ 16(DI), AX, R9
    ADOXQ AX, BP
    
    ADCXQ R9, SI
    MULXQ 24(DI), AX, R9
    ADOXQ AX, SI
    
    // add the last carries to R9
    MOVQ $0, DX
    ADCXQ DX, R9
    ADOXQ DX, R9
    
    // m := t[0]*q'[0] mod W
    MOVQ $0x0a117fffffffffff, DX
    MULXQ CX,R10, DX
    
    // clear the carry flags
    XORQ DX, DX
    
    // C,_ := t[0] + m*q[0]
    MOVQ $0x0a11800000000001, DX
    MULXQ R10, AX, DX
    ADCXQ CX ,AX
    MOVQ DX, CX
    
    // for j=1 to N-1
    //    (C,t[j-1]) := t[j] + m*q[j] + C
    MOVQ $0x59aa76fed0000001, DX
    ADCXQ  BX, CX
    MULXQ R10, AX, BX
    ADOXQ AX, CX
    MOVQ $0x60b44d1e5c37b001, DX
    ADCXQ  BP, BX
    MULXQ R10, AX, BP
    ADOXQ AX, BX
    MOVQ $0x12ab655e9a2ca556, DX
    ADCXQ  SI, BP
    MULXQ R10, AX, SI
    ADOXQ AX, BP
    MOVQ $0, AX
    ADCXQ AX, SI
    ADOXQ R9, SI
    
    // ---------------------------------------------------------------------------------------------
    // outter loop 2
    
    // clear up the carry flags
    XORQ R9 , R9
    
    // DX = y[2]
    MOVQ 16(R8), DX
    
    // for j=0 to N-1
    //    (A,t[j])  := t[j] + x[j]*y[i] + A
    
    MULXQ 0(DI), AX, R9
    ADOXQ AX, CX
    
    ADCXQ R9, BX
    MULXQ 8(DI), AX, R9
    ADOXQ AX, BX
    
    ADCXQ R9, BP
    MULXQ 16(DI), AX, R9
    ADOXQ AX, BP
    
    ADCXQ R9, SI
    MULXQ 24(DI), AX, R9
    ADOXQ AX, SI
    
    // add the last carries to R9
    MOVQ $0, DX
    ADCXQ DX, R9
    ADOXQ DX, R9
    
    // m := t[0]*q'[0] mod W
    MOVQ $0x0a117fffffffffff, DX
    MULXQ CX,R10, DX
    
    // clear the carry flags
    XORQ DX, DX
    
    // C,_ := t[0] + m*q[0]
    MOVQ $0x0a11800000000001, DX
    MULXQ R10, AX, DX
    ADCXQ CX ,AX
    MOVQ DX, CX
    
    // for j=1 to N-1
    //    (C,t[j-1]) := t[j] + m*q[j] + C
    MOVQ $0x59aa76fed0000001, DX
    ADCXQ  BX, CX
    MULXQ R10, AX, BX
    ADOXQ AX, CX
    MOVQ $0x60b44d1e5c37b001, DX
    ADCXQ  BP, BX
    MULXQ R10, AX, BP
    ADOXQ AX, BX
    MOVQ $0x12ab655e9a2ca556, DX
    ADCXQ  SI, BP
    MULXQ R10, AX, SI
    ADOXQ AX, BP
    MOVQ $0, AX
    ADCXQ AX, SI
    ADOXQ R9, SI
    
    // ---------------------------------------------------------------------------------------------
    // outter loop 3
    
    // clear up the carry flags
    XORQ R9 , R9
    
    // DX = y[3]
    MOVQ 24(R8), DX
    
    // for j=0 to N-1
    //    (A,t[j])  := t[j] + x[j]*y[i] + A
    
    MULXQ 0(DI), AX, R9
    ADOXQ AX, CX
    
    ADCXQ R9, BX
    MULXQ 8(DI), AX, R9
    ADOXQ AX, BX
    
    ADCXQ R9, BP
    MULXQ 16(DI), AX, R9
    ADOXQ AX, BP
    
    ADCXQ R9, SI
    MULXQ 24(DI), AX, R9
    ADOXQ AX, SI
    
    // add the last carries to R9
    MOVQ $0, DX
    ADCXQ DX, R9
    ADOXQ DX, R9
    
    // m := t[0]*q'[0] mod W
    MOVQ $0x0a117fffffffffff, DX
    MULXQ CX,R10, DX
    
    // clear the carry flags
    XORQ DX, DX
    
    // C,_ := t[0] + m*q[0]
    MOVQ $0x0a11800000000001, DX
    MULXQ R10, AX, DX
    ADCXQ CX ,AX
    MOVQ DX, CX
    
    // for j=1 to N-1
    //    (C,t[j-1]) := t[j] + m*q[j] + C
    MOVQ $0x59aa76fed0000001, DX
    ADCXQ  BX, CX
    MULXQ R10, AX, BX
    ADOXQ AX, CX
    MOVQ $0x60b44d1e5c37b001, DX
    ADCXQ  BP, BX
    MULXQ R10, AX, BP
    ADOXQ AX, BX
    MOVQ $0x12ab655e9a2ca556, DX
    ADCXQ  SI, BP
    MULXQ R10, AX, SI
    ADOXQ AX, BP
    MOVQ $0, AX
    ADCXQ AX, SI
    ADOXQ R9, SI
    
    reduce:
    // reduce, constant time version
    // first we copy registers storing t in a separate set of registers
    // as SUBQ modifies the 2nd operand
    MOVQ CX, DX
    MOVQ BX, R8
    MOVQ BP, R9
    MOVQ SI, R10
    MOVQ $0x0a11800000000001, R11
    SUBQ  R11, DX
    MOVQ $0x59aa76fed0000001, R11
    SBBQ  R11, R8
    MOVQ $0x60b44d1e5c37b001, R11
    SBBQ  R11, R9
    MOVQ $0x12ab655e9a2ca556, R11
    SBBQ  R11, R10
    JCS t_is_smaller // no borrow, we return t
    
    // borrow is set, we return u
    MOVQ DX, (DI)
    MOVQ R8, 8(DI)
    MOVQ R9, 16(DI)
    MOVQ R10, 24(DI)
    RET
    t_is_smaller:
    MOVQ CX, 0(DI)
    MOVQ BX, 8(DI)
    MOVQ BP, 16(DI)
    MOVQ SI, 24(DI)
    RET
    
    no_adx:
    
    // ---------------------------------------------------------------------------------------------
    // outter loop 0
    
    // (A,t[0]) := t[0] + x[0]*y[0]
    MOVQ (DI), AX // x[0]
    MOVQ 0(R8), R12
    MULQ R12 // x[0] * y[0]
    MOVQ DX, R9
    MOVQ AX, CX
    
    // m := t[0]*q'[0] mod W
    MOVQ $0x0a117fffffffffff, R10
    IMULQ CX , R10
    
    // C,_ := t[0] + m*q[0]
    MOVQ $0x0a11800000000001, AX
    MULQ R10
    ADDQ CX ,AX
    ADCQ $0, DX
    MOVQ  DX, R11
    
    // for j=1 to N-1
    //    (A,t[j])  := t[j] + x[j]*y[i] + A
    //    (C,t[j-1]) := t[j] + m*q[j] + C
    MOVQ 8(DI), AX
    MULQ R12 // x[1] * y[0]
    MOVQ R9, BX
    ADDQ AX, BX
    ADCQ $0, DX
    MOVQ DX, R9
    
    MOVQ $0x59aa76fed0000001, AX
    MULQ R10
    ADDQ  BX, R11
    ADCQ $0, DX
    ADDQ AX, R11
    ADCQ $0, DX
    
    MOVQ R11, CX
    MOVQ DX, R11
    MOVQ 16(DI), AX
    MULQ R12 // x[2] * y[0]
    MOVQ R9, BP
    ADDQ AX, BP
    ADCQ $0, DX
    MOVQ DX, R9
    
    MOVQ $0x60b44d1e5c37b001, AX
    MULQ R10
    ADDQ  BP, R11
    ADCQ $0, DX
    ADDQ AX, R11
    ADCQ $0, DX
    
    MOVQ R11, BX
    MOVQ DX, R11
    MOVQ 24(DI), AX
    MULQ R12 // x[3] * y[0]
    MOVQ R9, SI
    ADDQ AX, SI
    ADCQ $0, DX
    MOVQ DX, R9
    
    MOVQ $0x12ab655e9a2ca556, AX
    MULQ R10
    ADDQ  SI, R11
    ADCQ $0, DX
    ADDQ AX, R11
    ADCQ $0, DX
    
    MOVQ R11, BP
    MOVQ DX, R11
    
    ADDQ R11, R9
    MOVQ R9, SI
    
    // ---------------------------------------------------------------------------------------------
    // outter loop 1
    
    // (A,t[0]) := t[0] + x[0]*y[1]
    MOVQ (DI), AX // x[0]
    MOVQ 8(R8), R12
    MULQ R12 // x[0] * y[1]
    ADDQ AX, CX
    ADCQ $0, DX
    MOVQ DX, R9
    
    // m := t[0]*q'[0] mod W
    MOVQ $0x0a117fffffffffff, R10
    IMULQ CX , R10
    
    // C,_ := t[0] + m*q[0]
    MOVQ $0x0a11800000000001, AX
    MULQ R10
    ADDQ CX ,AX
    ADCQ $0, DX
    MOVQ  DX, R11
    
    // for j=1 to N-1
    //    (A,t[j])  := t[j] + x[j]*y[i] + A
    //    (C,t[j-1]) := t[j] + m*q[j] + C
    MOVQ 8(DI), AX
    MULQ R12 // x[1] * y[1]
    ADDQ R9, BX
    ADCQ $0, DX
    ADDQ AX, BX
    ADCQ $0, DX
    MOVQ DX, R9
    
    MOVQ $0x59aa76fed0000001, AX
    MULQ R10
    ADDQ  BX, R11
    ADCQ $0, DX
    ADDQ AX, R11
    ADCQ $0, DX
    
    MOVQ R11, CX
    MOVQ DX, R11
    MOVQ 16(DI), AX
    MULQ R12 // x[2] * y[1]
    ADDQ R9, BP
    ADCQ $0, DX
    ADDQ AX, BP
    ADCQ $0, DX
    MOVQ DX, R9
    
    MOVQ $0x60b44d1e5c37b001, AX
    MULQ R10
    ADDQ  BP, R11
    ADCQ $0, DX
    ADDQ AX, R11
    ADCQ $0, DX
    
    MOVQ R11, BX
    MOVQ DX, R11
    MOVQ 24(DI), AX
    MULQ R12 // x[3] * y[1]
    ADDQ R9, SI
    ADCQ $0, DX
    ADDQ AX, SI
    ADCQ $0, DX
    MOVQ DX, R9
    
    MOVQ $0x12ab655e9a2ca556, AX
    MULQ R10
    ADDQ  SI, R11
    ADCQ $0, DX
    ADDQ AX, R11
    ADCQ $0, DX
    
    MOVQ R11, BP
    MOVQ DX, R11
    
    ADDQ R11, R9
    MOVQ R9, SI
    
    // ---------------------------------------------------------------------------------------------
    // outter loop 2
    
    // (A,t[0]) := t[0] + x[0]*y[2]
    MOVQ (DI), AX // x[0]
    MOVQ 16(R8), R12
    MULQ R12 // x[0] * y[2]
    ADDQ AX, CX
    ADCQ $0, DX
    MOVQ DX, R9
    
    // m := t[0]*q'[0] mod W
    MOVQ $0x0a117fffffffffff, R10
    IMULQ CX , R10
    
    // C,_ := t[0] + m*q[0]
    MOVQ $0x0a11800000000001, AX
    MULQ R10
    ADDQ CX ,AX
    ADCQ $0, DX
    MOVQ  DX, R11
    
    // for j=1 to N-1
    //    (A,t[j])  := t[j] + x[j]*y[i] + A
    //    (C,t[j-1]) := t[j] + m*q[j] + C
    MOVQ 8(DI), AX
    MULQ R12 // x[1] * y[2]
    ADDQ R9, BX
    ADCQ $0, DX
    ADDQ AX, BX
    ADCQ $0, DX
    MOVQ DX, R9
    
    MOVQ $0x59aa76fed0000001, AX
    MULQ R10
    ADDQ  BX, R11
    ADCQ $0, DX
    ADDQ AX, R11
    ADCQ $0, DX
    
    MOVQ R11, CX
    MOVQ DX, R11
    MOVQ 16(DI), AX
    MULQ R12 // x[2] * y[2]
    ADDQ R9, BP
    ADCQ $0, DX
    ADDQ AX, BP
    ADCQ $0, DX
    MOVQ DX, R9
    
    MOVQ $0x60b44d1e5c37b001, AX
    MULQ R10
    ADDQ  BP, R11
    ADCQ $0, DX
    ADDQ AX, R11
    ADCQ $0, DX
    
    MOVQ R11, BX
    MOVQ DX, R11
    MOVQ 24(DI), AX
    MULQ R12 // x[3] * y[2]
    ADDQ R9, SI
    ADCQ $0, DX
    ADDQ AX, SI
    ADCQ $0, DX
    MOVQ DX, R9
    
    MOVQ $0x12ab655e9a2ca556, AX
    MULQ R10
    ADDQ  SI, R11
    ADCQ $0, DX
    ADDQ AX, R11
    ADCQ $0, DX
    
    MOVQ R11, BP
    MOVQ DX, R11
    
    ADDQ R11, R9
    MOVQ R9, SI
    
    // ---------------------------------------------------------------------------------------------
    // outter loop 3
    
    // (A,t[0]) := t[0] + x[0]*y[3]
    MOVQ (DI), AX // x[0]
    MOVQ 24(R8), R12
    MULQ R12 // x[0] * y[3]
    ADDQ AX, CX
    ADCQ $0, DX
    MOVQ DX, R9
    
    // m := t[0]*q'[0] mod W
    MOVQ $0x0a117fffffffffff, R10
    IMULQ CX , R10
    
    // C,_ := t[0] + m*q[0]
    MOVQ $0x0a11800000000001, AX
    MULQ R10
    ADDQ CX ,AX
    ADCQ $0, DX
    MOVQ  DX, R11
    
    // for j=1 to N-1
    //    (A,t[j])  := t[j] + x[j]*y[i] + A
    //    (C,t[j-1]) := t[j] + m*q[j] + C
    MOVQ 8(DI), AX
    MULQ R12 // x[1] * y[3]
    ADDQ R9, BX
    ADCQ $0, DX
    ADDQ AX, BX
    ADCQ $0, DX
    MOVQ DX, R9
    
    MOVQ $0x59aa76fed0000001, AX
    MULQ R10
    ADDQ  BX, R11
    ADCQ $0, DX
    ADDQ AX, R11
    ADCQ $0, DX
    
    MOVQ R11, CX
    MOVQ DX, R11
    MOVQ 16(DI), AX
    MULQ R12 // x[2] * y[3]
    ADDQ R9, BP
    ADCQ $0, DX
    ADDQ AX, BP
    ADCQ $0, DX
    MOVQ DX, R9
    
    MOVQ $0x60b44d1e5c37b001, AX
    MULQ R10
    ADDQ  BP, R11
    ADCQ $0, DX
    ADDQ AX, R11
    ADCQ $0, DX
    
    MOVQ R11, BX
    MOVQ DX, R11
    MOVQ 24(DI), AX
    MULQ R12 // x[3] * y[3]
    ADDQ R9, SI
    ADCQ $0, DX
    ADDQ AX, SI
    ADCQ $0, DX
    MOVQ DX, R9
    
    MOVQ $0x12ab655e9a2ca556, AX
    MULQ R10
    ADDQ  SI, R11
    ADCQ $0, DX
    ADDQ AX, R11
    ADCQ $0, DX
    
    MOVQ R11, BP
    MOVQ DX, R11
    
    ADDQ R11, R9
    MOVQ R9, SI
    
    JMP reduce
