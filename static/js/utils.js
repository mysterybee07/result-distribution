function fetchSemesters(fetchCoursesFlag = true) {
    const programId = document.getElementById('program_id').value;

    fetch(`/getfiltersemesters?program_id=${programId}`)
        .then(response => response.json())
        .then(data => {
            const semesterSelect = document.getElementById('semester_id');
            // semesterSelect.innerHTML = 'select';
            data.semesters.forEach(semester => {
                const option = document.createElement('option');
                option.value = semester.ID;
                option.textContent = semester.name;
                semesterSelect.appendChild(option);
            });

        })
        .catch(error => {
            console.error('Error fetching semesters:', error);
            alert('An error occurred while fetching semesters');
        });
}

function fetchCourses() {
    const programId = document.getElementById('program_id').value;
    const semesterId = document.getElementById('semester_id').value;

    fetch(`/getfiltercourses?program_id=${programId}&semester_id=${semesterId}`)
        .then(response => response.json())
        .then(data => {
            const courseSelect = document.getElementById('course_id');
            // courseSelect.innerHTML = '';
            data.courses.forEach(course => {
                const option = document.createElement('option');

                option.value = course.ID;
                option.textContent = course.name;
                courseSelect.appendChild(option);
            });
        })
        .catch(error => {
            console.error('Error fetching courses:', error);
            alert('An error occurred while fetching courses');
        });
}

document.addEventListener('DOMContentLoaded', () => {
    const form = document.getElementById('formScreen1'); // Replace with your form ID
    form.reset();
    // Fetch initial semesters and courses based on the selected program

    // fetchSemesters();

    // Event listeners for dynamic updates
    // document.getElementById('program_id').addEventListener('change', fetchSemesters);
    document.getElementById('semester_id').addEventListener('change', fetchCourses);

});