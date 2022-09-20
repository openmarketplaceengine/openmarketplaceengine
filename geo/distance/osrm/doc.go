/*
Package osrm provides functions over the Open Source Routing Machine project.

We use OSRM's Table service when we want to quickly determine the duration
and/or distance between two points (without traffic awareness). This could be
useful for saving money/requests to Google Maps API every time an anonymous end
user sees an estimate quote for a potential trip.
https://github.com/Project-OSRM/osrm-backend/blob/master/docs/http.md#table-service

We use OpenStreetMap's /reverse endpoint for reverse geocoding.
https://nominatim.org/release-docs/develop/api/Reverse/
*/
package osrm
