import React, { useState } from 'react'
import ListNotice from './ListNotice'
import { useQuery } from '@tanstack/react-query';
import { useData } from '../../context/DataContext';
import api from '../../api';

const AllNotice = () => {
    const { programs, batches } = useData();
    const [selectedProgram, setSelectedProgram] = useState();
    const [selectedBatch, setSelectedBatch] = useState();
    const fetchNotices = async () => {
        const apiUrl =
            selectedProgram && selectedBatch ? `/notice/by-program-and-batch?program_id=${selectedProgram}&batch_id=${selectedBatch}`
                : selectedProgram ? `/notice/by-program?program_id=${selectedProgram}`
                    : '/notice'; // Conditional URL based on programId
        const response = await api.get(apiUrl);
        console.log(response.data.notices);
        return Array.isArray(response.data.notices) ? response.data.notices : []; // Ensure it's always an array
    };

    // Using useQuery to fetch notices based on programId
    const { data, isLoading, isError, error } = useQuery({
        queryKey: ['notices', selectedProgram],  // Include programId in the query key
        queryFn: fetchNotices,
        enabled: true,  // This ensures the query runs immediately
    });

    if (isLoading) {
        return <p>Loading...</p>;
    }

    if (isError) {
        return <p>Error: {error.message}</p>;
    }
    return (
        <ListNotice
            data={data}
            programs={programs}
            batches={batches}
            selectedProgram={selectedProgram}
            setSelectedProgram={setSelectedProgram}
            setSelectedBatch={setSelectedBatch}
        />
    )
}

export default AllNotice