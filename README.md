# RC-Sender

This is just a little command application control 433mhz devices on Raspberry Pi. You don't need any external library like wiring-pi or pigpio. It's based on the [RC-Switch repository](https://github.com/sui77/rc-switch) for Arduino and ESP8266.

For rootless use (/dev/gpiomem) you need to add your user to the **GPIO** group.

## Build

```
./build/build-linux-your-architecture
```

## Usage

```
./rc-sender-linux-your-build <pin> <code> <length> <protocol> <pulse_length> <repeat_transmit>
```

| Argument        | Description                                                                 |
| --------------- | --------------------------------------------------------------------------- |
| pin             | GPIO Pin (bcm2835 pin, not physical pin - [layout](https://bit.ly/3uwtwzB)) |
| code            | your code                                                                   |
| length          | code length                                                                 |
| pulse_length    | custom pulse length (set to 0 if you want to use the default value)         |
| repeat_transmit | number of repeats (set to 0 if you want to use the default value)           |
