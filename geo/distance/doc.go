/*
Package distance provides functions that abstract over various 3rd-party
Distance Matrix APIs, such as OSRM, Google Maps, GraphHopper, MapBox, etc.

Typically, in Distance Matrix APIs, you pass in a list of N origins and M
destinations, and the result is an NxM matrix of elements that contain
information about distance/duration between every point combination.
*/
package distance
