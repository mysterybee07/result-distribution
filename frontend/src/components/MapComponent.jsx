import React from 'react';
import { MapContainer, TileLayer, Marker, Popup } from 'react-leaflet';
import 'leaflet/dist/leaflet.css';

const MapComponent = ({ colleges }) => {
    console.log("🚀 ~ MapComponent ~ college:", colleges)
    // Define the initial position of the map
    const center = [27.6663423, 85.3330053]; // Latitude and Longitude (e.g., London)

    // Define marker positions
    const markers = colleges.map((college) => ({
        position: [parseFloat(college.latitude), parseFloat(college.longitude)],
        content: college.college_name,
      }));

    return (
        <MapContainer center={center} zoom={15} style={{ width: "100%", height: "50vh" }}>
            {/* Add a tile layer (map background) */}
            <TileLayer
                url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
            />

            {/* Add markers */}
            {markers.map((marker, index) => (
                <Marker key={index} position={marker.position}>
                    <Popup>{marker.content}</Popup>
                </Marker>
            ))}
        </MapContainer>
    );
};

export default MapComponent;