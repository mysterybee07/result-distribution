import React, { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import api from "../../api";
import {
  flexRender,
  useReactTable,
  getCoreRowModel,
  getPaginationRowModel,
  getSortedRowModel,
} from "@tanstack/react-table";
import { Button } from "@/components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { FaEdit, FaTrash } from "react-icons/fa";
import { useNavigate } from "react-router-dom";
import { Input } from "@/components/ui/input";
import { UpdateCenter } from "../../components/UpdateCenter";

export default function ListCollege() {
  const navigate = useNavigate();
  const [sorting, setSorting] = useState([]);
  const [filter, setFilter] = useState("");
  const [selectedColleges, setSelectedColleges] = useState([]);

  console.log("🚀 ~ ListCollege ~ selectedColleges:", selectedColleges)

  const fetchColleges = async () => {
    const response = await api.get("/college");
    console.log("🚀 ~ fetchColleges ~ response:", response.data.center);
    return response.data.colleges;
  };

  const { data: college = [], isLoading, error } = useQuery({
    queryKey: ["colleges"],
    queryFn: fetchColleges,
  });
  console.log("🚀 ~ ListCollege ~ college:", college)

  const handleCheckboxChange = (id) => {
    setSelectedColleges((prev) =>
      prev.includes(id)
        ? prev.filter((collegeId) => collegeId !== id)
        : [...prev, id]
    );
  };

  const columns = [
    {
      id: "select",
      header: ({ table }) => (
        <input
          type="checkbox"
          onChange={(e) =>
            e.target.checked
              ? setSelectedColleges(college.map((item) => item.college_name))
              : setSelectedColleges([])
          }
          checked={selectedColleges.length === college.length && college.length > 0}
        />
      ),
      cell: ({ row }) => (
        <input
          type="checkbox"
          checked={selectedColleges.includes(row.original.college_name)}
          onChange={() => handleCheckboxChange(row.original.college_name)}
          disabled={row.original.is_center} // Disable if is_center is true
        />
      ),
    },
    {
      accessorKey: "ID",
      header: "S.N",
      cell: ({ row }) => row.index + 1,
    },
    {
      accessorKey: "college_name",
      header: "Name",
    },
    {
      accessorKey: "address",
      header: "Address",
    },
    {
      id: "actions",
      header: "Action",
      cell: ({ row }) => {
        const data = row.original;
        return (
          <div className="flex items-center gap-4">
            <FaEdit
              className="text-blue-600 cursor-pointer"
              onClick={() => navigate(`/admin/students/edit/${data.ID}`)}
            />
            <FaTrash
              className="text-red-600 cursor-pointer"
              onClick={() => navigate(`/admin/students/${data.ID}`)}
            />
            <UpdateCenter center={data.is_center} capacity={data.capacity} id={data.ID}/>
          </div>
        );
      },
    },
  ];

  const table = useReactTable({
    data: college,
    columns,
    state: {
      sorting,
      globalFilter: filter,
    },
    onSortingChange: setSorting,
    getCoreRowModel: getCoreRowModel(),
    // getPaginationRowModel: getPaginationRowModel(),
    getSortedRowModel: getSortedRowModel(),
  });

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;

  return (
    <div className="w-full">
      <div className="flex text-left items-center justify-between py-4">
        <Input
          placeholder="Search colleges..."
          value={filter}
          onChange={(e) => setFilter(e.target.value)}
          className="max-w-sm"
        />
        <div className="flex gap-2">
          {selectedColleges.length > 0 &&
            <Button
              onClick={() => navigate(`/admin/center/create?colleges=${encodeURIComponent(JSON.stringify(selectedColleges))}`)}
              size="sm"
            >
              Update as center
            </Button>

          }
          <Button onClick={() => navigate("/admin/college/create")} size="sm">
            Add College
          </Button>
        </div>
      </div>
      <div className="rounded-md border">
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow key={headerGroup.id}>
                {headerGroup.headers.map((header) => (
                  <TableHead key={header.id}>
                    {header.isPlaceholder
                      ? null
                      : flexRender(header.column.columnDef.header, header.getContext())}
                  </TableHead>
                ))}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {table.getRowModel().rows.length ? (
              table.getRowModel().rows.map((row) => (
                <TableRow key={row.id}>
                  {row.getVisibleCells().map((cell) => (
                    <TableCell key={cell.id} className="text-left">
                      {flexRender(cell.column.columnDef.cell, cell.getContext())}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell colSpan={columns.length} className="text-left">
                  No results found.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
      <div className="flex items-center justify-between py-4">
        <span className="text-sm">
          {table.getRowModel().rows.length} of {college.length} row(s)
        </span>
        <div className="space-x-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => table.previousPage()}
            disabled={!table.getCanPreviousPage()}
          >
            Previous
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={() => table.nextPage()}
            disabled={!table.getCanNextPage()}
          >
            Next
          </Button>
        </div>
      </div>
    </div>
  );
}


