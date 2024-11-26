import React, { useState } from 'react'
import {
    Card,
    CardContent,
    CardDescription,
    CardFooter,
    CardHeader,
    CardTitle,
} from '../../components/ui/card'
import api from '../../api';
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { useNavigate } from 'react-router-dom';
import { FaEdit } from "react-icons/fa";
import { MdDelete } from "react-icons/md";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import {
    Select,
    SelectContent,
    SelectGroup,
    SelectItem,
    SelectLabel,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select"
import { set } from 'react-hook-form';



const ListNotice = ({
    data,
    batches,
    programs,
    selectedProgram,
    setSelectedProgram,
    setSelectedBatch,
}) => {
    const navigate = useNavigate();

    const handleEdit = (id) => {
        console.log(id);
        navigate(`/admin/notice/edit/${id}`);
        // TODO: Redirect to edit page
    };
    const handleDelete = async (id) => {
        console.log(id);
        const res = await api.delete(`/notice/delete/${id}`);
    };

    function getSemesterSuffix(semester) {
        if (semester === 1) return 'st';
        if (semester === 2) return 'nd';
        if (semester === 3) return 'rd';
        return 'th'; // For 4th, 5th, 6th, etc.
    }

    const [p, setP] = useState(null);
    const [b, setB] = useState(null);


    return (
        <div>
            <div className='flex flex-row justify-between mb-2'>
                <p className="text-dark text-start font-bold text-xl">Recent Notices</p>
                <Button onClick={() => navigate('/admin/notice/create')}>Create Notice</Button>
            </div>
            <div className='space-x-4 mb-4'>
                <div className='flex flex-row space-x-4'>
                    {/* <p>Filer:</p> */}
                    <Select
                        id="program-select"
                        className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                        onValueChange={(value) => setP(value)} // On change, call filterByProgram with selected program ID
                    >
                        <SelectTrigger>
                            <SelectValue placeholder="Select program" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectGroup>
                                {/* <SelectItem value="">Select Program</SelectItem> Default SelectItem */}
                                {programs?.map((program) => (
                                    <SelectItem key={program.ID} value={program.ID}>
                                        {program.program_name}
                                    </SelectItem>
                                ))}
                            </SelectGroup>
                        </SelectContent>
                    </Select>
                    {p && (
                        <Select
                            id="program-select"
                            className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                            onValueChange={(value) => setB(value)} // On change, call filterByProgram with selected program ID
                        >
                            <SelectTrigger>
                                <SelectValue placeholder="All Batch" />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectGroup>
                                    {/* <SelectItem value="">Select Program</SelectItem> Default SelectItem */}
                                    {batches?.map((batch) => (
                                        <SelectItem key={batch.ID} value={batch.ID}>
                                            {batch.batch}
                                        </SelectItem>
                                    ))}
                                </SelectGroup>
                            </SelectContent>
                        </Select>
                    )}
                    <Button variant="default" onClick={() => {
                        console.log(p, b);
                        setSelectedProgram(p);
                        setSelectedBatch(b);
                    }} >Filter</Button>
                    <Button variant="secondary" onClick={() => {
                        setP('');
                        setB('');
                        setSelectedProgram('');
                        setSelectedBatch('');
                    }} >Clear</Button>
                </div>
            </div>
            <div>
                {data?.map((notice, index) => (
                    <Card
                        key={notice.ID}
                        className={`relative ${index !== data.length - 1 ? "mb-4" : ""} group`} // Add 'group' for hover targeting
                    >
                        {/* Edit and Delete Button Container */}
                        <div
                            className="absolute top-2 right-2 flex space-x-2 opacity-0 group-hover:opacity-100 transition-opacity duration-200"
                        >
                            <Button
                                variant="outline"
                                size="icon"
                                onClick={() => handleEdit(notice.ID)}
                            >
                                <FaEdit />
                            </Button>
                            <Button
                                variant="destructive"
                                size="icon"
                                onClick={() => handleDelete(notice.ID)}
                            >
                                <MdDelete />
                            </Button>
                        </div>

                        <CardHeader>
                            <CardTitle>{notice.Title}</CardTitle>
                            <CardDescription>{(notice?.Created_at)?.split("T")[0]}</CardDescription>
                            <CardDescription className="space-x-2">
                                <Badge variant="secondary">{notice.Batch}</Badge>
                                <Badge variant="secondary">{notice.Program}</Badge>
                                <Badge variant="secondary">
                                    {notice.Semester}
                                    <sup>{getSemesterSuffix(notice.Semester)}</sup> &nbsp; Semester                                </Badge>
                            </CardDescription>
                        </CardHeader>

                        <CardContent>
                            <p>{notice.Description}</p>
                        </CardContent>

                        {notice.FilePath && (
                            <CardFooter>
                                <Button
                                    onClick={() =>
                                        window.open("/1732362345505245562-Coverletter.pdf", "_blank")
                                    }
                                >
                                    Open File
                                </Button>
                            </CardFooter>
                        )}
                    </Card>
                ))}
            </div>
        </div>
    )
}

export default ListNotice