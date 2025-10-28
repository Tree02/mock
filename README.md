
# MockLogin

Mock preparado para simular puestos operativos de rp5, cargados mediante archivo config y manejando generación de JWT para retornar a quien los consulte.

Mock de emulación Moca Operacional proporciona token a los distintos servicios que lo requieran **_RapiGo_**

Resuelve momentáneamente requerimiento del área de desarrollo con el fin de seguir flujos de programa dependientes de Mantra.

### Funciones
Al momento de ejecución del programa se leeran los puestos introducidos en el archivo de configuración, y se cargarán en un mapa de puestos.

El proyecto cuenta con los **endpoints**

* /token/get
* /token/refresh
* /token/bye

Casos de **respuesta**
```
OK
{
    "details": "ok",
    "data":{
        "jwt":"",
        "start": ""
    } 
}
```
```
ERROR
{
    "details": "failed",
    "error": ""
}
```

Cada uno con una función distinta.

En el caso de /token/get, de ser la primer consulta, generará el token y lo devolverá como respuesta dentro de data.jwt, y el timestamp de cuándo se generó, dentro del campo data.start

En el caso de /token/refresh, devolverá el token nuevo generado del puesto en el que se esté consultando, propiamente en data.jwt, y el timestamp en data.start

En el caso de /token/bye se abandonará el puesto que se esté consultando y la iteración pasará al siguiente, sumandole +1 al indice

El programa quedará indefinidamente corriendo hasta que se corte el flujo del mismo, volviendo al índice 1 al momento de terminar de recorrer el mapa de puestos.

Si se quiere ejecutar en background 
```
nohup ./mocklogin & 
```

### Finalización
Una vez cortemos el flujo del programa, este borrará las variables creadas de JWT, START, etc. Cortará el listener en la ip y puerto especificados.
