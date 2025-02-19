import React, { useMemo } from 'react';
import {
    ResponsiveContainer,
    LineChart,
    Line,
    CartesianGrid,
    XAxis,
    YAxis,
    Tooltip,
    Legend
} from 'recharts';

const calculateSemesterGPA = (marks) => {
    console.log("🚀 ~ calculateSemesterGPA ~ marks:", marks)
    // Group marks by semester
    const semesterMarks = marks.reduce((acc, mark) => {
        const sem = mark.semester_id;
        if (!acc[sem]) {
            acc[sem] = [];
        }
        acc[sem].push(mark);
        return acc;
    }, {});

    // Calculate GPA for each semester
    return Object.entries(semesterMarks).map(([semester, marks]) => {
        const totalPoints = marks.reduce((sum, mark) => {
            // Convert percentage to GPA (example scale)
            const percentage = (mark.total_marks / (mark.Course.semester_total_marks +
                mark.Course.practical_total_marks +
                mark.Course.assistant_total_marks)) * 100;
            let gradePoint;
            if (percentage >= 90) gradePoint = 4.0;
            else if (percentage >= 80) gradePoint = 3.7;
            else if (percentage >= 70) gradePoint = 3.3;
            else if (percentage >= 60) gradePoint = 3.0;
            else if (percentage >= 50) gradePoint = 2.7;
            else gradePoint = 2.0;

            return sum + gradePoint;
        }, 0);

        const averageGPA = (totalPoints / marks.length).toFixed(2);

        return {
            semester: `Semester ${semester}`,
            gpa: parseFloat(averageGPA)
        };
    }).sort((a, b) => a.semester.localeCompare(b.semester));
};

const StudentGPAChart = ({ marks }) => {
    console.log("🚀 ~ marksGPAChart ~ marks:", marks)
    const semesterResults = useMemo(() => {
        if (!marks?.length) return [];
        return calculateSemesterGPA(marks);
    }, [marks]);

    return (
        <div className="w-full h-[400px]">
            <ResponsiveContainer width="100%" height="100%">
                <LineChart
                    data={semesterResults}
                    margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
                >
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="semester" />
                    <YAxis domain={[0, 4]} />
                    <Tooltip />
                    <Legend />
                    <Line
                        type="monotone"
                        dataKey="gpa"
                        stroke="#4F46E5"
                        activeDot={{ r: 8 }}
                        strokeWidth={2}
                    />
                </LineChart>
            </ResponsiveContainer>
        </div>
    );
};

export default StudentGPAChart;