# streamdeck-plugin-hlae
 [WIP] StreamDeck plugin to communicate with CSGO via HLAE's command(mirv_pgl)

## How to use
Execute following command on CS:GO client with HLAE :   
```mirv_pgl url "ws://localhost:65535/std"```  
then   
```mirv_pgl start```   
to start communicating with StreamDeck plugin.

### CAUTION
If you streamdeck app is down, which means HLAE Server is not running, you should execute ```mirv_pgl stop```  to avoid mirvpgl polling. Otherwise your CS:GO client will start stuttering.