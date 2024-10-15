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
import { useQuery } from "@tanstack/react-query";
import api from "../api";

const invoices = [
    {
        invoice: "INV001",
        paymentStatus: "Paid",
        totalAmount: "$250.00",
        paymentMethod: "Credit Card",
    },
    {
        invoice: "INV002",
        paymentStatus: "Pending",
        totalAmount: "$150.00",
        paymentMethod: "PayPal",
    },
    {
        invoice: "INV003",
        paymentStatus: "Unpaid",
        totalAmount: "$350.00",
        paymentMethod: "Bank Transfer",
    },
    {
        invoice: "INV004",
        paymentStatus: "Paid",
        totalAmount: "$450.00",
        paymentMethod: "Credit Card",
    },
    {
        invoice: "INV005",
        paymentStatus: "Paid",
        totalAmount: "$550.00",
        paymentMethod: "PayPal",
    },
    {
        invoice: "INV006",
        paymentStatus: "Pending",
        totalAmount: "$200.00",
        paymentMethod: "Bank Transfer",
    },
    {
        invoice: "INV007",
        paymentStatus: "Unpaid",
        totalAmount: "$300.00",
        paymentMethod: "Credit Card",
    },
]

export default function StudentTable() {
    const { data: students, isLoading, error } = useQuery({
        queryKey: ['students'],
        queryFn: async () => {
            const response = await api.get("/students");
            return response.data.students;
        },
    });
    console.log(students);

    if (isLoading) return <div>Loading...</div>;
    if (error) return <div>Error: {error.message}</div>;
    return (
        <>
            <Table>
                <TableCaption className="font-bold text-xl">Student Table</TableCaption>
                <TableHeader>
                    <TableRow>
                        <TableHead>Name</TableHead>
                        <TableHead>Symbol Number</TableHead>
                        <TableHead>Reg Number</TableHead>
                        <TableHead>Batch</TableHead>
                        <TableHead>Program</TableHead>
                        <TableHead>Semester</TableHead>
                        <TableHead>Status</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {Array.isArray(students) && students.map((student) => (
                        <TableRow key={student.ID}>
                            <TableCell>{student.fullname}</TableCell>
                            <TableCell>{student.symbol_number}</TableCell>
                            <TableCell>{student.registration_number}</TableCell>
                            <TableCell>{student.Batch.batch}</TableCell>
                            <TableCell>{student.Program.program_name}</TableCell>
                            <TableCell>{student.current_semester}</TableCell>
                            <TableCell>{student.status}</TableCell>
                        </TableRow>
                    ))}
                </TableBody>
                {/* <TableFooter>
                    <TableRow>
                        <TableCell colSpan={3} className="text-start">Total</TableCell>
                        <TableCell className="text-right">$2,500.00</TableCell>
                    </TableRow>
                </TableFooter> */}
            </Table>
        </>
    )
}
