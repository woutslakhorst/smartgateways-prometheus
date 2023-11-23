# smartgateways-prometheus
Docker image for transforming smartgateways-kamstir metrics to prometheus metrics

The program is hardwired to fetch data from http://connectix_kamstir.local:82/kamst-ir/api/read

Metrics are exposed on port: 8080 path: /metrics

The following metrics are exported:

- kamstir_gj_total from "heat_energy"
- kamstir_temp1_c_current from "temp1"
- kamstir_temp2_c_current from "temp2"
- kamstir_tempdiff_c_current from "tempdiff"
- kamstir_flow_m3h_current from "flow"
- kamstir_volume_m3_total from "volume"

The smartgateways module is queried when the API is called.
The smartgateways module only updates it's internal data every 10 minutes.