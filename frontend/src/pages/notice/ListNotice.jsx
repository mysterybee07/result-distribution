import React from 'react'
import {
    Card,
    CardContent,
    CardDescription,
    CardFooter,
    CardHeader,
    CardTitle,
} from '../../components/ui/card'
import { useQuery } from "@tanstack/react-query";
import api from '../../api';
import { Badge } from "@/components/ui/badge"


const ListNotice = () => {
    const fetchNotices = async () => {
        const response = await api.get('/notice');
        console.log(response.data.notices);
        return response.data.notices;
    };

    const { data, isLoading, isError, error } = useQuery({
        queryKey: ['notices'],
        queryFn: fetchNotices,
    });

    if (isLoading) {
        return <p>Loading...</p>;
    }

    if (isError) {
        return <p>Error: {error.message}</p>;
    }

    return (
        <div>
            {data.map((notice) => (
                <Card key={notice.ID}>
                    <CardHeader>
                        <CardTitle>{notice.title}</CardTitle>
                        <CardDescription>
                            <Badge variant="outline">{notice.batch_id}</Badge>
                            <Badge variant="outline">{notice.program_id.program_name}</Badge>
                            <Badge variant="outline">{notice.semester_id}</Badge>
                        </CardDescription>

                        <CardDescription>{(notice.CreatedAt).split("T")[0]}</CardDescription>
                    </CardHeader>
                    <CardContent>
                        <p>{notice.description}</p>
                    </CardContent>
                    {/* <CardFooter>
                    <p>Card Footer</p>
                </CardFooter> */}
                </Card>
                // <li key={notice.id}>{notice.title}</li> // Assuming notices have an `id` and `title`
            ))}



        </div>
    )
}

export default ListNotice