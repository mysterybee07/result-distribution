import React from 'react'
import CourseForm from '../../forms/CourseForm'
import { useParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';

const EditCourse = () => {
    const { id } = useParams();
    console.log("🚀 ~ EditCourse ~ id:", id)
  return (
    <div>
        <CourseForm />
    </div>
  )
}

export default EditCourse