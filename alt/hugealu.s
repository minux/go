TEXT	main(SB),512|7,$0
	MOVUB	0xf00ddd(R1), R2
	MOVD	-16529(R2), R3
	ADD	$0xf00abcd, R1
	AND	$0xf00abcd, R1, R2
	ADD	$0x1f00abcd, R3
	MOVD	R4, 0xf00ddd(R5)
	MOVD	R6, -16521(R8)
	MOVUB	ZR, -16529(R8)
	MOVUB	R1, 0xf00ddd(R8)
	RET
