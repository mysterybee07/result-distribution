import { useQuery } from '@tanstack/react-query'
import React from 'react'

const ListCollege = () => {
  const fetchColleges = async () => {
    const response = await api.get("/colleges");
    return response.data;
};
  const { data: college =[], isLoading, error } = useQuery({
    queryKey: ['colleges'],
    queryFn: fetchColleges,
  })
  if (isLoading) return <div>Loading...</div>;
    if (error) return <div>Error: {error.message}</div>;
  return (
    <div>ListCollege</div>
  )
}

export default ListCollege