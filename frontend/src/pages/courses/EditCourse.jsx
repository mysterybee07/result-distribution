import React, { useEffect } from 'react'
import CourseForm from '../../forms/CourseForm'
import { useParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import api from '../../api';

const EditCourse = () => {
  const { id } = useParams();

  const fetchCourseById = async () => {
    const response = await api.get(`/courses/${id}`);
    return response.data.course;
  }

  const { data: course = {}, isLoading, isError, error } = useQuery({
    queryKey: ["course", id],
    queryFn: fetchCourseById,
    enabled: !!id,
  });
  console.log("🚀 ~ EditCourse ~ course:", course)
  return ( 
    <div>
      <CourseForm data= {course}/>
    </div>
  )
}

export default EditCourse