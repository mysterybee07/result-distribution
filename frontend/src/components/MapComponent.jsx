// import { MapContainer, TileLayer, Marker, Popup } from "react-leaflet";
// import "leaflet/dist/leaflet.css";

// const MapComponent = ({ lat, lng }) => {
//     return (
//         <MapContainer center={[lat, lng]} zoom={13} style={{ height: "400px", width: "100%" }}>
//             <TileLayer
//                 url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
//             />
//             <Marker position={[lat, lng]}>
//                 <Popup>Location: {lat}, {lng}</Popup>
//             </Marker>
//         </MapContainer>
//     );
// };

// export default MapComponent;

import { GoogleMap, LoadScript, Marker } from "@react-google-maps/api";

const MapComponent = ({ lat, lng }) => {
    const mapContainerStyle = { height: "400px", width: "100%" };
    const center = { lat, lng };

    return (
        <LoadScript googleMapsApiKey="YOUR_GOOGLE_MAPS_API_KEY">
            <GoogleMap mapContainerStyle={mapContainerStyle} center={center} zoom={13}>
                <Marker position={center} />
            </GoogleMap>
        </LoadScript>
    );
};

export default MapComponent;
