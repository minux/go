TEXT	main(SB),512|7,$0
	MOVD	$32(BSP), R4
	MOVD	32(BFP), R5
	MOVD	BFP, R1
	MOVD	R2, BSP
	MOVD	R3, 16(BSP)
	ADD     $16, BSP
	RET
