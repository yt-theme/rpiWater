#!/usr/bin/python3
#_*_ coding:utf-8 _*_
#coding=utf-8

import sys
import getopt

try:
    import RPi.GPIO as GPIO
except:
    print("import gpio err")

args = sys.argv

# 参数 1,pin号(num)     2,输入输出模式(IN/OUT)     3,读取值/高低电平(READ/HIGH/LOW)
# 参数 1,init(INIT) 2,x 3,x 进行gpio初始化操作
# 参数 1,cleanup(CLEANUP) 2,x 3,x 进行gpio恢复状态操作

if args[1] == "INIT":
    pass
    
if args[1] == "CLEANUP":
    GPIO.cleanup()

# #########################################################################

elif args[1] != "INIT":

    # set GPIO mode
    GPIO.setmode(GPIO.BOARD)
    # close warning
    GPIO.setwarnings(False)
    
    # set IO mode
    if args[2] == "OUT":
        GPIO.setup(int(args[1]), GPIO.OUT)
        
    if args[2] == "IN":
        GPIO.setup(int(args[1]), GPIO.IN)
        
    
    # set IO
    if args[3] == "READ":
        readRet = GPIO.input(int(args[1]))
        print(readRet)
    if args[3] == "HIGH":
        GPIO.output(int(args[1]), GPIO.HIGH)
        print(args[1], "HIGH =>")
    if args[3] == "LOW":
        GPIO.output(int(args[1]), GPIO.LOW)
        print(args[1], "LOW =>")
# while 1:
#     pass