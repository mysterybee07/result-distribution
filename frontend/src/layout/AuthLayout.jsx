import {
    SidebarInset,
    SidebarProvider,
    SidebarTrigger,
} from "@/components/ui/sidebar"
import {
    Breadcrumb,
    BreadcrumbItem,
    BreadcrumbLink,
    BreadcrumbList,
    BreadcrumbPage,
    BreadcrumbSeparator,
} from "@/components/ui/breadcrumb"
import { Separator } from "@/components/ui/separator"
import { AppSidebar } from "../components/AppSidebar";
import { Toaster } from '@/components/ui/toaster';
import { Outlet } from 'react-router-dom';
import Footer from "../components/Footer";
import { useLocation, Link } from "react-router-dom";

const AdminLayout = () => {
    const location = useLocation(); // Get the current location
    const { pathname } = location;

    // Helper to split the path and exclude the `/admin` prefix
    const generateBreadcrumbs = () => {
        const segments = pathname
            .split("/")
            .filter((segment) => segment && segment !== "admin"); // Exclude 'admin' and empty segments

        return segments.map((segment, index) => {
            // Construct the path for each breadcrumb
            const url = `/${["admin", ...segments.slice(0, index + 1)].join("/")}`; // Re-add `/admin` for correct routing
            return {
                name: formatBreadcrumbName(segment),
                url
            };
        });
    };

    // Helper to format breadcrumb names (e.g., replace hyphens, capitalize words)
    const formatBreadcrumbName = (name) => {
        return name
            .replace(/-/g, " ") // Replace hyphens with spaces
            .replace(/\b\w/g, (char) => char.toUpperCase()); // Capitalize first letter of each word
    };

    const breadcrumbs = generateBreadcrumbs();
    return (
        <>
            <SidebarProvider>
                <AppSidebar />
                <SidebarInset>
                    <header className="flex sticky top-0 bg-dark h-12 shrink-0 items-center gap-2 border-b px-4">
                        <SidebarTrigger className="-ml-1" />
                        <Separator orientation="vertical" className="mr-2 h-4" />
                        <Breadcrumb>
                            <BreadcrumbList>
                                {/* Map over the breadcrumb array */}
                                {breadcrumbs.map((breadcrumb, index) => (
                                    <BreadcrumbItem key={breadcrumb.url}>
                                        {index !== breadcrumbs.length - 1 ? (
                                            <>
                                                <BreadcrumbLink asChild>
                                                    <Link to={breadcrumb.url}>{breadcrumb.name}</Link>
                                                </BreadcrumbLink>
                                                <BreadcrumbSeparator />
                                            </>
                                        ) : (
                                            <BreadcrumbPage>{breadcrumb.name}</BreadcrumbPage>
                                        )}
                                    </BreadcrumbItem>
                                ))}
                            </BreadcrumbList>
                        </Breadcrumb>
                    </header>
                    <div className="flex flex-1 flex-col gap-1 py-1 px-8">
                        {/* <div className="flex-grow"> */}
                            <Toaster />
                            <Outlet />
                        {/* </div> */}
                        <Footer />
                    </div>
                </SidebarInset>
            </SidebarProvider>

            {/* <AdminNavbar /> */}



        </>
    );
};

export default AdminLayout;