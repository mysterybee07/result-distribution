import React from 'react'
import { Button } from '../components/ui/button'
import { NavLink, Link, useNavigate } from 'react-router-dom'
import { CgProfile } from "react-icons/cg";
import { IoIosLogOut, IoIosNotifications } from "react-icons/io";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { useAuth } from "../context/AuthContext"

const AdminNavbar = () => {
    const navigate = useNavigate()
    const { isAuthenticated, logout } = useAuth();
    return (
        <nav className="fixed inset-x-0 top-0 z-50 bg-white shadow-sm dark:bg-gray-950/90">
            <div className="w-full max-w-7xl mx-auto px-4">
                <div className="flex justify-between h-14 items-center">
                    <Link href="#" className="flex items-center" prefetch={false}>
                        {/* <MountainIcon className="h-6 w-6" /> */}
                        <img src='../../public/logo.png' alt='logo' className='w-8 h-8 cursor-pointer rounded-full' />
                        <span className="ml-2 text-lg font-semibold">Result-e</span>
                    </Link>
                    <nav className="hidden md:flex gap-16">
                        <NavLink
                            to="/admin"
                            className={({ isActive }) =>
                                isActive
                                    ? "font-medium flex items-center text-sm transition-colors underline text-blue-600" // Active link styling
                                    : "font-medium flex items-center text-sm transition-colors hover:underline"
                            }
                            prefetch={false}
                        >
                            Home
                        </NavLink>
                        <NavLink
                            to="/admin/exam"
                            className={({ isActive }) =>
                                isActive
                                    ? "font-medium flex items-center text-sm transition-colors underline text-blue-600" // Active link styling
                                    : "font-medium flex items-center text-sm transition-colors hover:underline"
                            } prefetch={false}
                        >
                            Exam
                        </NavLink>
                        <NavLink
                            to="/admin/result"
                            className={({ isActive }) =>
                                isActive
                                    ? "font-medium flex items-center text-sm transition-colors underline text-blue-600" // Active link styling
                                    : "font-medium flex items-center text-sm transition-colors hover:underline"
                            } prefetch={false}
                        >
                            Result
                        </NavLink>
                    </nav>
                    <div className="flex items-center gap-4">
                        {isAuthenticated ?
                            (
                                <>
                                    <DropdownMenu>
                                        <DropdownMenuTrigger>
                                            <IoIosNotifications className='cursor-pointer w-6 h-6' />
                                        </DropdownMenuTrigger>
                                        <DropdownMenuContent>
                                            <DropdownMenuLabel className="text-orange-600">2 new notification</DropdownMenuLabel>
                                            <DropdownMenuSeparator />
                                            <DropdownMenuItem>
                                                Study for exam..
                                            </DropdownMenuItem>
                                            <DropdownMenuItem>
                                                Final term exam in 7 days
                                            </DropdownMenuItem>
                                        </DropdownMenuContent>
                                    </DropdownMenu>

                                    <DropdownMenu>
                                        <DropdownMenuTrigger>
                                            <img src='../../public/neko.jpg' alt='profile' className='w-8 h-8 cursor-pointer rounded-full' />
                                        </DropdownMenuTrigger>
                                        <DropdownMenuContent>
                                            <DropdownMenuLabel>Kiran Shrestha</DropdownMenuLabel>
                                            <DropdownMenuSeparator />
                                            <DropdownMenuItem>
                                                <CgProfile className='cursor-pointer mr-1' />Profile
                                            </DropdownMenuItem>
                                            <DropdownMenuItem onClick={logout}>
                                                <IoIosLogOut className='cursor-pointer mr-1' />
                                                Logout
                                            </DropdownMenuItem>
                                        </DropdownMenuContent>
                                    </DropdownMenu>
                                </>
                            ) :
                            (
                                <>
                                    <Button variant="outline" size="sm" onClick={() => navigate('/login')}>
                                        Sign in
                                    </Button>
                                    <Button size="sm" onClick={() => navigate('/register')}>Sign up</Button>
                                </>
                            )}


                    </div>
                </div>
            </div>
        </nav>

    )
}

export default AdminNavbar