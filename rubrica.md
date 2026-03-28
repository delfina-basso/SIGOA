# Rúbrica para el Trabajo Práctico Grupal - SIGOA

Puntaje Total: **100 puntos** (80 puntos para la implementación general y 20 puntos para los aspectos generales)

## Módulo 1: Check-in de Pasajeros (20 puntos)

1. **Representación de Vuelos y Pasajeros (Informe - 2 puntos)**:

   Descripción clara y justificada de las estructuras de datos (listas, colas, etc.) utilizadas para representar vuelos, pasajeros y reservas. Se explica cómo se relacionan estas estructuras.

2. **Implementación de la Búsqueda de Pasajeros (Implementación - 3 puntos)**:

   Implementación funcional y eficiente de la búsqueda por DNI/Pasaporte, recuperando toda la información requerida correctamente.

3. **Simulación de Llegada y Cola de Check-in (Implementación - 5 puntos)**:

   Simulación lógica y configurable de la llegada de pasajeros y la gestión de una cola única con múltiples mostradores.

4. **Proceso de Check-in (Implementación - 5 puntos)**:

   Implementación completa del proceso, incluyendo validación de pasajeros, despacho de equipaje (simulado), emisión de tarjeta de embarque y gestión de la lista de pasajeros del vuelo.

5. **Gestión de la Lista de Espera (Implementación - 5 puntos)**:

   Implementación robusta de la lista de espera, con lógica correcta de asignación por prioridad (platino > oro > plata > no frecuente) y orden de llegada. Actualización precisa de la nómina.

## Módulo 2: Gestión de Cargas (20 puntos)

1. **Representación de Cargas y Aeronaves (Informe - 2 puntos)**:

   Descripción clara y justificada de las estructuras de datos utilizadas para representar cargas y aeronaves, incluyendo sus capacidades.

2. **Lectura de Datos de Carga (Implementación - 3 puntos)**:

   Lectura correcta y eficiente de la información de los paquetes de carga desde el archivo, incluyendo destino, peso y volumen.

3. **Implementación de Restricciones de Capacidad (Implementación - 5 puntos)**:

   Implementación precisa de la verificación de las restricciones de peso y volumen de las aeronaves al intentar asignar carga.

4. **Algoritmo de Asignación de Carga (Implementación - 10 puntos)**:

   Implementación eficiente y lógica de un algoritmo (o heurística) para maximizar el uso del espacio y respetar las restricciones de peso. Justificación clara de la elección del algoritmo y su funcionamiento.

## Módulo 3: Embarque de Pasajeros (20 puntos)

1. **Representación de Zonas y Asientos (Informe - 2 puntos)**:

   Descripción clara y justificada de cómo se representan las zonas de embarque y cómo se utiliza la información del archivo `configuracion_asientos.txt`.

2. **Lógica de Embarque Prioritario (Implementación - 2 puntos)**:
   Implementación correcta y eficiente del embarque prioritario para las categorías Platino, Oro y Plata.
3. **Lógica de Embarque por Zonas (Implementación - 3 puntos)**:

   Implementación precisa y lógica del embarque secuencial por zonas, utilizando correctamente la información del archivo `configuracion_asientos.txt` para determinar la zona de cada pasajero.

4. **Gestión de Pasajeros No Presentados (Implementación - 5 puntos)**:

   Implementación correcta de la lógica para identificar y dar de baja a los pasajeros que no se presentan a la hora límite de embarque.

5. **Generación del Archivo de Vuelo Comprimido (Implementación - 8 puntos)**:

   Implementación correcta y eficiente de la recopilación de toda la información relevante del vuelo y la compresión utilizando el algoritmo de Huffman. Generación del archivo comprimido funcional y correcto.

## Módulo 4: Despacho de Vuelos (20 puntos)

1. **Implementación de la Cola de Prioridad (Implementación - 5 puntos)**:

   Implementación correcta y eficiente de la cola de prioridad para gestionar el despacho de vuelos, considerando la hora de salida programada y el estado de preparación del vuelo.

2. **Lectura de Datos de Edificios (Implementación - 5 puntos)**:

   Lectura correcta y eficiente de la información de los edificios y obstáculos desde el archivo `edificios.txt`.

3. **Algoritmo de Cálculo de la Línea del Horizonte (Implementación - 10 puntos)**:

   Implementación eficiente y correcta del algoritmo para calcular la línea del horizonte a partir de los datos de los edificios. Generación precisa y correcta del archivo de salida.

## Aspectos Generales (20 puntos)

1. **Organización y Calidad del Código (Implementación - 10 puntos)**:

   Código bien estructurado, modularizado, comentado de forma clara y consistente, siguiendo buenas prácticas de programación en Go. Uso adecuado de nombres de variables y funciones.

2. **Calidad del Informe Escrito (General - 4 puntos)**:

   Informe bien estructurado siguiendo la plantilla, claro, conciso, completo y con un lenguaje técnico adecuado y sin errores gramaticales. Coherencia entre las descripciones de los diferentes módulos.

3. **Calidad de la Presentación Oral (General - 6 puntos)**:

   Presentación clara, bien organizada, con un flujo lógico y un buen uso del tiempo. Todos los miembros del grupo demuestran un conocimiento sólido del trabajo en su totalidad y responden de manera precisa y bien articulada a las preguntas generales.

## Distribución de Puntos

| **Categoría**         | **Subcategoría**                               | **Puntos** |
| --------------------- | ---------------------------------------------- | ---------- |
| Check in              | Vuelos y Pasajeros                             | 2          |
|                       | Búsqueda de Pasajeros                          | 3          |
|                       | Simulación de Llegada y Cola de Check-in       | 5          |
|                       | Proceso de Check-in                            | 5          |
|                       | Gestión de la Lista de Espera                  | 5          |
| Gestión de Cargas     | Representación de Cargas y Aeronaves           | 2          |
|                       | Lectura de Datos de Carga                      | 3          |
|                       | Implementación de Restricciones de Capacidad   | 5          |
|                       | Algoritmo de Asignación de Carga               | 10         |
| Embarque de Pasajeros | Representación de Zonas y Asientos             | 2          |
|                       | Lógica de Embarque Prioritario                 | 2          |
|                       | Lógica de Embarque por Zonas                   | 3          |
|                       | Gestión de Pasajeros No Presentados            | 5          |
|                       | Generación del Archivo de Vuelo Comprimido     | 8          |
| Despacho de Vuelos    | Implementación de la Cola de Prioridad         | 5          |
|                       | Lectura de Datos de Edificios                  | 5          |
|                       | Algoritmo de Cálculo de la Línea del Horizonte | 10         |
| Aspectos Generales    | Organización y Calidad del Código              | 10         |
|                       | Calidad del Informe Escrito                    | 4          |
|                       | Calidad de la Presentación Oral                | 6          |
| **Total**             |                                                | **100**    |
