import {
    Card,
    CardContent,
    CardDescription,
    CardFooter,
    CardHeader,
    CardTitle,
} from "@/components/ui/card"
import {
    Tabs,
    TabsContent,
    TabsList,
    TabsTrigger,
} from "@/components/ui/tabs"

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
import ProgramForm from "../forms/ProgramForm"
import BatchForm from "../forms/BatchForm"
import { useQuery } from "@tanstack/react-query"
import api from "../api"

const programs = [
    {
        sn: "1",
        program: "CSIT",
        student: "100+"
    },
    {
        sn: "2",
        program: "BBA",
        student: "100+"
    },
    {
        sn: "3",
        program: "BIM",
        student: "100+"
    },
    {
        sn: "4",
        program: "BBS",
        student: "100+"
    },
]

const batches = [
    {
        sn: "1",
        batch: "2024",
        student: "100+"
    },
    {
        sn: "2",
        batch: "2023",
        student: "100+"
    },
    {
        sn: "3",
        batch: "2022",
        student: "100+"
    },
    {
        sn: "4",
        batch: "2021",
        student: "100+"
    },
    {
        sn: "5",
        batch: "2020",
        student: "100+"
    },
]

export default function DashboardAsideCard() {
    const { data} = useQuery({
        queryKey: "programs",
        queryFn: async () => {
            const response = await api.get("/program");
            return response.data;
        },
    });
    console.log("data:",data);

    return (
        <Tabs defaultValue="program" className="w-[400px]">
            <TabsList className="grid w-full grid-cols-2">
                <TabsTrigger value="program">Programs</TabsTrigger>
                <TabsTrigger value="batches">Batches</TabsTrigger>
            </TabsList>
            <TabsContent value="program">
                <Card>
                    <CardHeader>
                        <CardTitle>Programs</CardTitle>
                        <CardDescription>
                            These are the listed programs in the server.
                        </CardDescription>
                    </CardHeader>
                    <CardContent className="space-y-2">
                        <Table>
                            <TableHeader>
                                <TableRow>
                                    <TableHead className="w-[50px] text-center">S.N</TableHead>
                                    <TableHead className="w-[150px] text-center">Programs</TableHead>
                                    <TableHead className="w-[150px] text-center">Students</TableHead>
                                </TableRow>
                            </TableHeader>
                            <TableBody>
                                {programs.map((program) => (
                                    <TableRow key={program.program}>
                                        <TableCell className="font-medium text-center">{program.sn}</TableCell>
                                        <TableCell className="font-medium text-center">{program.program}</TableCell>
                                        <TableCell className="font-medium text-center">{program.student}</TableCell>
                                    </TableRow>
                                ))}
                            </TableBody>
                        </Table>
                    </CardContent>
                    <CardFooter>
                        <ProgramForm />
                    </CardFooter>
                </Card>
            </TabsContent>

            <TabsContent value="batches">
                <Card>
                    <CardHeader>
                        <CardTitle>Batches</CardTitle>
                        <CardDescription>
                            These are the listed batches in the server.
                        </CardDescription>
                    </CardHeader>
                    <CardContent className="space-y-2">
                        <Table>
                            <TableHeader>
                                <TableRow>
                                    <TableHead className="w-[50px] text-center">S.N</TableHead>
                                    <TableHead className="w-[150px] text-center">Programs</TableHead>
                                    <TableHead className="w-[150px] text-center">Students</TableHead>
                                </TableRow>
                            </TableHeader>
                            <TableBody>
                                {batches.map((batch) => (
                                    <TableRow key={batch.batch}>
                                        <TableCell className="font-medium text-center">{batch.sn}</TableCell>
                                        <TableCell className="font-medium text-center">{batch.batch}</TableCell>
                                        <TableCell className="font-medium text-center">{batch.student}</TableCell>
                                    </TableRow>
                                ))}
                            </TableBody>
                        </Table>
                    </CardContent>
                    <CardFooter>
                        <BatchForm />
                    </CardFooter>
                </Card>
            </TabsContent>
        </Tabs>
    )
}
