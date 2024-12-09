import { useQuery } from '@tanstack/react-query'
import React from 'react'
import api from '../../api';
import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableFooter,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { FaEdit, FaTrash } from 'react-icons/fa';
import { Button } from '../../components/ui/button';
import { useNavigate } from 'react-router-dom';

const ListCollege = () => {
  const navigate = useNavigate();
  const fetchColleges = async () => {
    const response = await api.get("/college");
    console.log("🚀 ~ fetchColleges ~ response:", response.data.center)
    return response.data.center;
  };
  const { data: college = [], isLoading, error } = useQuery({
    queryKey: ['colleges'],
    queryFn: fetchColleges,
  })
  // console.log(college)
  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;
  return (
    <div>
      <div className='flex justify-between items-center'>
        <TableCaption className='text-center'>List of colleges</TableCaption>
        < Button onClick={() => navigate('/admin/college/create')} size="sm" > Add College</Button>
      </div>
      <Table className="text-left">
        <TableHeader>
          <TableRow>
            <TableHead>
              S.N
            </TableHead>
            <TableHead>
              Name
            </TableHead>
            <TableHead>
              Address
            </TableHead>
            <TableHead>
              Action
            </TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {college.map((data, index) => (
            <TableRow>
              <TableCell>{index + 1}</TableCell>
              <TableCell>{data.college_name}</TableCell>
              <TableCell>{data.address}</TableCell>
              <TableCell className="flex items-center gap-4">
                {/* todo change after route has been added */}
                <FaEdit
                  className="text-blue-600 cursor-pointer"
                  onClick={() => navigate(`/admin/students/edit/${data.ID}`)}
                />
                <FaTrash
                  onClick={() => navigate(`/admin/students/${data.ID}`)}
                  className="text-red-600"
                />
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  )
}

export default ListCollege