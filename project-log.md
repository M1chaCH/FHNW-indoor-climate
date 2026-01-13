# Project Log
This documents contains a simple documentation of what I did and what I plan to do. I'll update it after every working session.

## Session 01
- Unterstand Assignment
- Workout plan
- Getting started with micro:bit
- Get started with reading sensor data
- Write a go lang app that reads the data from the serial port an stores it in a csv file.
- use python and the provided template to draw plots

### Getting started with micro:bit
Reminder: Connect the USB Cable to the micro:bit controller and not to the grove board. When connected to the grove board, you cannot upload micro:bit programs. (There is no disk mounted to the computer.)

### Getting started with reading sensor data
(AI and Google were used to find answers for some of the questions)  
**What does the number in the screen command do?**  
`screen /dev/tty.usbmodem102 115200`  

This Number is the Baud Rate. This number is defined by the writing device. The reader must match this rate, otherwise the data will be read wrongly.  

What is a Baud Rate?: https://riverdi.com/blog/understanding-baud-rate-a-comprehensive-guide  
It defines the number of signal changes per second. The higher the Baud Rate the more data gets through per second. 

Bit rate is something similar. Usually a Baud Rate of 1 results in 1bps (bit per second). But there are different scenarios where you can "encode" more than one bit into one signal change resulting in higher throughput. 



**What is the unit of the data measured?**  
C02:  
- Provided in ppm (parts per Million) 
    - 10'000ppm = 1% of "Air" is CO2
- Range: 400ppm - 10'000ppm  
- Accuracy: +- 3%
- sensor measures the value using NDIR (Non-Dispersive Infrared)

Humidity:  
Measured in %. 
100% means the Air has reached the maximum amount of water vapor possible at the current temperature. When the temperature changes then the maximum amount also changes. 
If the amount were to increase above 100% then the air becomes water. 

Temperature:  
You can choose between Fahrenheit and Celsius.

**Why the I2C port?**  
Because this is an "intelligent" port. Over this port actual "structured" data is sent. The other ports just receive an analog or digital signal. (raw voltage)

### Use python
I have never properly used python. This is why I've consulted Gemini for help here. I did not copy code directly but I certainly received a lot of help from the LLM.


## Session 02
