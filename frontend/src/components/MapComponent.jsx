import React from 'react';
import { MapContainer, TileLayer, Marker, Popup } from 'react-leaflet';
import 'leaflet/dist/leaflet.css';

const MapComponent = ({ college }) => {
    console.log("🚀 ~ MapComponent ~ college:", college)
    // Define the initial position of the map
    const center = [27.6663423, 85.3330053]; // Latitude and Longitude (e.g., London)

    // Define marker positions
    const markers = [
        { position: [27.6663423, 85.3330053], content: "Marker 1" },
        { position: [51.51, -0.1], content: "Marker 2" },
        { position: [51.49, -0.08], content: "Marker 3" },
    ];

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