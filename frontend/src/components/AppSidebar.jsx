import * as React from "react"
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible"
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarRail,
  SidebarFooter,
} from "@/components/ui/sidebar"
import { ChevronRightIcon } from "@radix-ui/react-icons"
import { Link, useLocation } from "react-router-dom"
// import { NavUser } from "@/components/nav-user"
import { NavUser} from "./ui/nav-user"
import { useAuth } from "../context/AuthContext"


// This is sample data.
const data = {
  navMain: [
    {
      title: "Course",
      url: "/admin/courses",
      items: [
        {
          title: "View Course",
          url: "/admin/courses",
        },
        {
          title: "Add Course",
          url: "/admin/courses/create",
        },
      ],
    },
    {
      title: "College",
      url: "/admin/college",
      items: [
        {
          title: "View College",
          url: "/admin/college",
        },
        {
          title: "Add College",
          url: "/admin/college/create",
        },
        {
          title: "List Center",
          url: "/admin/center",
        },
      ],
    },
    {
      title: "Student",
      url: "/admin/students",
      items: [
        {
          title: "View Student",
          url: "/admin/students",
        },
        {
          title: "Add Student",
          url: "/admin/students/create",
        },
      ],
    },
    {
      title: "Result",
      url: "/admin/result",
      items: [
        {
          title: "View Result",
          url: "/admin/result",
        },
        {
          title: "Add Result",
          url: "/admin/result",
        },
      ],
    },
    {
      title: "Exam",
      url: "/admin/exam",
      items: [
        {
          title: "View Exam",
          url: "/admin/exam",
        },
        {
          title: "Routine",
          url: "/admin/exam/routine",
        },
        {
          title: "Exam Schedule",
          url: "/admin/exam/create",
        },
      ],
    },
    {
      title: "Notice",
      url: "/admin/notice",
      items: [
        {
          title: "View Notice",
          url: "/admin/notice",
        },
        {
          title: "Add Notice",
          url: "/admin/notice",
        },
      ],
    },
  ],
}

export function AppSidebar({ ...props }) {
  const location = useLocation(); // Get the current location
  const currentPath = location.pathname; // Extract the pathname
  // console.log("🚀 ~ AppSidebar ~ currentPath:", currentPath)
  const { isAuthenticated, logout, userData } = useAuth();
  // console.log("🚀 ~ AppSidebar ~ userData:", userData)
  const user = {
    // TODO: change this no name later
    name: userData?.role,
    email: userData?.email,
    avatar: userData?.avatar,
  }

  return (
    <Sidebar {...props}>
      <SidebarHeader>
        Result-e
      </SidebarHeader>

      <SidebarContent className="gap-0">
        <SidebarGroup>
          <Link
            to="/dashboard"
          >
            <SidebarMenuButton className={currentPath === "/dashboard" ? " font-semibold bg-sidebar-accent" : "font-semibold text-sidebar-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground"}
            >

              Home
            </SidebarMenuButton>
          </Link>
        </SidebarGroup>
        {/* We create a collapsible SidebarGroup for each parent. */}
        {data.navMain.map((item) => (
          <Collapsible
            key={item.title}
            title={item.title}
            defaultOpen
            className="group/collapsible"
          >
            <SidebarGroup>
              <SidebarGroupLabel
                asChild
                className="group/label text-sm font-semibold text-sidebar-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground"
              >
                <CollapsibleTrigger>
                  {item.title}{""}
                  <ChevronRightIcon className="ml-auto transition-transform group-data-[state=open]/collapsible:rotate-90" />
                </CollapsibleTrigger>
              </SidebarGroupLabel>

              <CollapsibleContent>
                <SidebarGroupContent>
                  <SidebarMenu>
                    {item.items.map((subItem) => {
                      // Determine if the current route is active
                      const isActive = currentPath === subItem.url;

                      return (
                        <SidebarMenuItem key={subItem.title}>
                          <SidebarMenuButton asChild isActive={isActive} >
                            <Link
                              to={subItem.url}
                              className={isActive ? "font-extrabold bg-slate-900" : ""}
                            >
                              {subItem.title}
                            </Link>
                          </SidebarMenuButton>
                        </SidebarMenuItem>
                      );
                    })}
                  </SidebarMenu>
                </SidebarGroupContent>
              </CollapsibleContent>
            </SidebarGroup>
          </Collapsible>
        ))}
      </SidebarContent>
      <SidebarFooter>
        <NavUser user={user} logout={logout}/>
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
}

