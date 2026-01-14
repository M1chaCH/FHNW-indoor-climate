# FHNW (ipro) Indoor Climate
FHNW (ipro) is a mandatory individual software project, worth 6 ETCS.

## Overview
In [this project](http://www.tamberg.org/fhnw/2025/hs/IproIndoorClimate.pdf) I will use CO2 sensors to measure indoor climate.

## Levels
To keep you motivated, this project template is split into levels.

- [x] Level 0: [Getting started](level-0/README.md)
- [x] Level 1: [Logging sensor data](level-1/README.md)
- [x] (SKIPPED) Level 2: [Analyzing your data](level-2/README.md)
- [x] Level 3: [Monitoring remotely](level-3/README.md)
- [ ] Level 4: [Scaling up and out](level-4/README.md)

## Results
Each level results in a working prototype, built from building blocks.

### Level 1
My prototype after this level is a working sensor that sends data over serial, a script written in go that reads the serial port and stores the data in csv files and a python script that reads the files and plots a chart using the data.
In regard to my project goal (getting to know MQTT) this concludes Level 1 for me. I could certainly improve this prototype either by drawing plots using live data and/or by storing the data in a proper database.

### Level 2
This level was skipped.

### Level 3
<kbd><img src="l3_kibana_dashboard.png" height="320"/></kbd>

My prototype from the third level is a working sensor that sends json data every 5 seconds to an online web api. This service then enriches the data and sends it to an Elasticsearch DB. The data can be read and evaluated in Kibana. The [Kibana Dashboard can be accessed here](https://kibana.michu-tech.com/s/indoor-climate/app/dashboards#/view/9df80984-140a-48b2-b364-bd9b4bb9c807?_g=(filters:!(),refreshInterval:(pause:!t,value:60000),time:(from:now-60m,to:now))). Use the User `ipro` to login. Reach out to me for the password.

This concludes level 3 for me. In the next level, I'll try to improve some of the shortcomings of this prototype.
## Goals
I am used to writing backends and frontends to process data, but working with sensors is entirely new to me.
So my goal for this project is to understand how I can work with simple sensors. How do I have to receive and process the data inorder to do something meaningful with it.  
Ultimately, I want to understand and use MQTT. 

Steps: 
- [x] Get to know micro:bit
- [x] read and process sensor data over a serial connection
- [x] read and process sensor data over a WiFi connection
- [ ] implement MQTT service to receive data from n sensors and process their data
- [x] implement backend that sends structured sensor data to elastic
- [x] create a simple kibana dashboard
- [ ] (optional) handwrite my own web frontend for data visualization

## Language
I am trying to experiment with languages. I'll use whatever programming language I think is most suitable for the task at hand.

## "AI" tools
No LLM code is copied into this repository directly. But depending on the task, I'll be consulting LLMs more. Similar to how I used to use Google to learn something, I now use LLMs to learn about programming patterns. In the project-log.md I will keep track of the tasks where I heavily relied on LLM input.

## Support
Contact thomas.amberg@fhnw.ch to get an MS Teams invite.

> Note: Work in progress. Interested? Contact thomas.amberg@fhnw.ch

## License
Unless noted otherwise.

* Source code examples in this repository are declared Public Domain [CC0 1.0](https://creativecommons.org/publicdomain/zero/1.0/)
* Content by [A. Kennel](https://www.fhnw.ch/de/personen/andrea-kennel), [G. Deck](https://www.fhnw.ch/en/people/klaus-georg-deck), [T. Amberg](https://www.fhnw.ch/en/people/thomas-amberg), FHNW is licensed under [CC BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/)

Publishing your own code?
[MIT License](https://choosealicense.com/licenses/mit/)
