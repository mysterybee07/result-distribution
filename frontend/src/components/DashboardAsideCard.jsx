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
import { useState } from "react"
import { Button } from "@/components/ui/button"
import { useData } from "../context/DataContext"

export default function DashboardAsideCard() {
    const {
        programs,
        loadingPrograms,
        errorPrograms,
        batches,
        loadingBatches,
        errorBatches,
    } = useData();

    if (loadingPrograms && loadingBatches) return <div>Loading...</div>;

    if (errorPrograms && errorBatches) return <div>Error: {errorBatches.message || errorPrograms.message}</div>;

    return (
        <Tabs defaultValue="program" className="w-[400px]">
            <TabsList className="grid w-full grid-cols-3">
                <TabsTrigger value="program">Programs</TabsTrigger>
                <TabsTrigger value="batches">Batches</TabsTrigger>
                <TabsTrigger value="semester">Semester</TabsTrigger>
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
                                    <TableHead className="w-[50px] text-center">Action</TableHead>
                                </TableRow>
                            </TableHeader>
                            <TableBody>
                                {Array.isArray(programs) && programs.map((program, index) => (
                                    <TableRow key={program.ID}>
                                        <TableCell className="font-medium text-center">{index + 1}</TableCell>
                                        <TableCell className="font-medium text-center">{program.program_name}</TableCell>
                                        <TableCell className="font-medium text-center">
                                            {/* <Button size="sm" variant="outline" onClick={<ProgramForm />} >
                                                Edit
                                            </Button> */}
                                            <ProgramForm program={program} />
                                        </TableCell>
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
                                    <TableHead className="w-[150px] text-center">Batches</TableHead>
                                    <TableHead className="w-[50px] text-center">Action</TableHead>
                                </TableRow>
                            </TableHeader>
                            <TableBody>
                                {Array.isArray(batches) && batches.map((batch, index) => (
                                    <TableRow key={batch.batch}>
                                        <TableCell className="font-medium text-center">{index + 1}</TableCell>
                                        <TableCell className="font-medium text-center">{batch.batch}</TableCell>
                                        <TableCell className="font-medium text-center">
                                            <BatchForm batch={batch} />
                                        </TableCell>

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

            <TabsContent value="semester">
                <Card>
                    <CardHeader>
                        <CardTitle>Semesters</CardTitle>
                        <CardDescription>
                            These are the listed semesters in the server.
                        </CardDescription>
                    </CardHeader>
                    <CardContent className="space-y-2">
                        <Table>
                            <TableHeader>
                                <TableRow>
                                    <TableHead className="w-[50px] text-center">S.N</TableHead>
                                    <TableHead className="w-[150px] text-center">Semesters</TableHead>
                                    {/* <TableHead className="w-[150px] text-center">Students</TableHead> */}
                                </TableRow>
                            </TableHeader>
                            <TableBody>
                                {Array.isArray(batches) && batches.map((batch, index) => (
                                    <TableRow key={batch.batch}>
                                        <TableCell className="font-medium text-center">{index + 1}</TableCell>
                                        <TableCell className="font-medium text-center">{batch.batch}</TableCell>
                                        {/* <TableCell className="font-medium text-center">{batch.student}</TableCell> */}
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
