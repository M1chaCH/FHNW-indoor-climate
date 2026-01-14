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
- Understand the FeatherS3
- send sensor data over serial
- send sensor data over WiFi
- create service that receives the sent data over WiFi
- send sensor data to elastic
- create small kibana dashboard


### Getting started with FeatherS3
Trying to understand the code from the documentation:
Sources: 
- https://docs.circuitpython.org/projects/scd30/en/latest/
- https://gemini.google.com/share/0ad50dbfe1d2

```python
import time
import board
import adafruit_scd30

# What does the following comment mean?
#   "SCD-30 has temperamental I2C with clock stretching, datasheet recommends starting at 50KHz"
# "temperamental I2C" -> the sensor returns data using the I2C protocol and has its own serial clock.
# "clock stretching" -> Master (FeatherS3) requests data at a fixed frequency. But the Slave (sensor) needs time to process the analogue data actually read by the sensor. This means the slave cannot provide accurate data at a fixed frequency. To solve this the slave does clock stretching (letting the master wait until he has finished processing the sensor data)
#
# Why did this line not work? What is the difference to the line I've written
#   "i2c = busio.I2C(board.SCL, board.SDA,frequency=50000)"
# The default line did not work because it was missing the import statement of "busio"...
#
# board.SCL points to the pin with the serial clock
# board.SDA points to the pin with the serial data
# busio.I2C configures the two GPIO (not specific purpose pins, can be programmed to anything) pins (SCL, SDA) to use the I2C protocol
# frequency=50000: do not change this. 50kHz is the most reliable. Increasing it might work but will probably result in wrong data from the CO2 sensor.
i2c = board.I2C() # this probably just a wrapper to the busio call.

# initialized an object with the high level api to access the sensor data.
scd = adafruit_scd30.SCD30(i2c)

while True:
    # What is meant by this?
    #   "since the measurement interval is long (2+ seconds) we check for new data before reading the values, to ensure current readings."
    # Due to physical constraints the CO2 value can only be measured approximately every 2 seconds. If we poll more frequently we won't see a real change in the value (maybe some jitter due to inaccurate floating point operations)
    # this interval could be changed here: scd.measurement_interval
    # do not assign a value lower than 2 because this will result in inaccurate data, but the value can always be increased.
    if scd.data_available:
        print("Data Available!")
        print("CO2:", scd.CO2, "PPM")
        print("Temperature:", scd.temperature, "degrees C")
        print("Humidity:", scd.relative_humidity, "%%rH")
        print("")
        print("Waiting for new data...")
        print("")

    time.sleep(0.5)
```

### Connecting the ESP to my mac using WiFi
https://learn.adafruit.com/networking-in-circuitpython/networking-with-the-wifi-module

### Level 3 Prototype
I'll create a C# service that exposes an HTTP endpoint, which will receive data from the ESP and send it to elastic. 
Kibana will then visualize this data.
The C# service lives in a docker container exposed to the internet. 
The ELK Stack is also already setup in my local network and is exposed to the internet.
The ESP with have authenticate itself by the backend using an API key.

#### Configuring Elastic indicies properly
I found, that elastic supports timeseries data storage. This is optimal for my usecase, but I need to read into this technology and figure out how I can apply it for my usecase.

https://www.elastic.co/docs/manage-data/data-store/data-streams/set-up-tsds

First, I have to create an index template. This is what this template does: 
- Doesn't create indices but a data stream
    - a data stream is a set of automatically managed indices. The application can just send data to this data stream (looks like an index from the api), but elastic will manage large amounts of data much more efficiently
- enables Timeseries Mode
- allows for automatically created indices
- configures the timeseries data structure: 
    - defines the index.routing_path
        - in my case: 
            - sensor_type
            - device_id
        - these are the fields that uniquely identify a timeseries (kind of like the group by in an SQL select statement)

Then I created a component template. I there, i defined the mapping of the fields to their type:
- device_id: keyword, timeseries dimension
- sensor_type: keyword, timeseries dimension
- co2: double, timeseries gauge metric
- temperature: double, timeseries gauge metric
- humidity: double, timeseries gauge metric
- the field @timestamp is automatically added as a date.

#### Kibana Dashboard
Further Idea: 
Visualize the angle of increase per 15 minutes or so (how quickly is the air quality getting worse, at what time of the day is the worst, is it getting exponentially worse or just linearly)


## Session 3
This session should result in a complete product. Here is a small brainstorm of what could be done:

```mermaid
mindmap
  root(Dynamic Sensor Hub)
    (mqtt)
      support n amounts of sensors
      support more types of sensors, maybe even any
        microphone? | measure volume
        moisture sensor | measure plant soil moisture
        human presence sensor$
        dust sensor
      sensors can register on their own
    (configuration)
      sensors can be configured from a remote source
      sensors push their available options
      EVT: remote source code bearbeiten | es kann irgend ein program auf den microchip eingespielt werden
    (reliability)
      microchips that loose connection to web service don't loose data accuracy
      microchips send metrics data
      battery power?
    (UI)
      a page analogue Kibana dashboard | kibana dashboards can be embedded
      a page where you can configure the sensors
      use websockets for "live web UI"
    (backend)
      cqrs architecure?
      .net or go?
      proper observability | send metrics and logs to elastic using open telemetry
    (elastic)
      dynamic mappings? | if i want to support any sensor

```