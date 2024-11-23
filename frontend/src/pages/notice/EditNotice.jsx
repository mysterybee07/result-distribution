import React from 'react'
import NoticeForm from '../../forms/NoticeForm'
import { useQuery } from '@tanstack/react-query';
import api from '../../api';

const EditNotice = () => {
    const url = window.location.pathname;
    const id = url.substring(url.lastIndexOf('/') + 1);
    const getNoticeById = async () => {
        const response = await api.get('/notice/by-id/' + id);
        const notice = response.data.notice;
        return notice;
    };

    const { data, isLoading, isError, error } = useQuery({
        queryKey: ['notices'],
        queryFn: getNoticeById,
    });
    if (isLoading) {
        return <p>Loading...</p>;
    }

    if (isError) {
        return <p>Error: {error.message}</p>;
    }
    
    return (
        <NoticeForm
            mode='edit'
            notice={data}
        />
    )
}

export default EditNotice