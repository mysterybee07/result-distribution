import * as React from "react"
import { Minus, Plus } from "lucide-react"
import { Button } from "@/components/ui/button"
import {
    Drawer,
    DrawerClose,
    DrawerContent,
    DrawerDescription,
    DrawerFooter,
    DrawerHeader,
    DrawerTitle,
    DrawerTrigger,
} from "@/components/ui/drawer"

export function UpdateCenter({ center }) {
    const [goal, setGoal] = React.useState(1)

    function onClick(adjustment) {
        setGoal(goal + adjustment)
    }

    const handleClick = () => {
        const is_center = true;
        const college_name = '';
        const capacity = '';

    }
    return (
        <Drawer>
            <DrawerTrigger asChild>
                <Button
                    variant="outline"
                    disabled={!center}
                    style={{
                        backgroundColor: center ? "green" : "transparent",
                        color: center ? "white" : "black",
                        borderColor: center ? "green" : "gray",
                    }}
                >
                    Center
                </Button>

            </DrawerTrigger>
            <DrawerContent>
                <div className="mx-auto w-full max-w-sm">
                    <DrawerHeader>
                        <DrawerTitle>Update Center</DrawerTitle>
                        <DrawerDescription>Add the number of students.</DrawerDescription>
                    </DrawerHeader>
                    <div className="p-4 pb-0">
                        <div className="flex items-center justify-center space-x-2">
                            <Button
                                variant="outline"
                                size="icon"
                                className="h-8 w-8 shrink-0 rounded-full"
                                onClick={() => onClick(-1)}
                                disabled={goal <= 0}
                            >
                                <Minus />
                                <span className="sr-only">Decrease</span>
                            </Button>
                            <div className="flex-1 text-center">
                                <div className="text-7xl font-bold tracking-tighter">
                                    <input
                                        name="goal"
                                        // type="number"
                                        value={goal === 0 ? "" : goal} // Show empty input when goal is 0
                                        min={1}
                                        max={1000}
                                        onChange={(event) => {
                                            const value = event.target.value;
                                            if (value === "") {
                                                setGoal(0); // Temporarily set goal to 0 for an empty input
                                            } else {
                                                const newValue = Math.max(1, Math.min(1000, Number(value))); // Enforce min/max constraints
                                                setGoal(newValue);
                                            }
                                        }}
                                        className="w-full text-center text-7xl font-bold tracking-tighter bg-transparent border-none outline-none"
                                    />


                                </div>
                                <div className="text-[0.70rem] uppercase text-muted-foreground">
                                    Students
                                </div>
                            </div>

                            <Button
                                variant="outline"
                                size="icon"
                                className="h-8 w-8 shrink-0 rounded-full"
                                onClick={() => onClick(1)}
                                disabled={goal >= 1000}
                            >
                                <Plus />
                                <span className="sr-only">Increase</span>
                            </Button>
                        </div>
                    </div>
                    <DrawerFooter>
                        <Button
                            onClick={handleClick}
                        >
                            Submit
                        </Button>
                        <DrawerClose asChild>
                            <Button variant="outline">Cancel</Button>
                        </DrawerClose>
                    </DrawerFooter>
                </div>
            </DrawerContent>
        </Drawer>
    )
}
