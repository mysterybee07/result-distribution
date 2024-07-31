let shouldFetchCourses = false; // Initialize with false if you don't want to fetch courses initially

function fetchSemesters() {
    const programId = document.getElementById('program_id').value;
    console.log('Fetching semesters for program ID:', programId);

    fetch(`/getfiltersemesters?program_id=${programId}`)
        .then(response => {
            console.log('Semesters response status:', response.status);
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            const semesterSelect = document.getElementById('semester_id');
            semesterSelect.innerHTML = ''; // Clear existing options

            // Add the default option
            const defaultOption = document.createElement('option');
            defaultOption.value = '';
            defaultOption.textContent = 'Select Semester';
            semesterSelect.appendChild(defaultOption);

            // Add new options
            data.semesters.forEach(semester => {
                const option = document.createElement('option');
                option.value = semester.ID;
                option.textContent = semester.name;
                semesterSelect.appendChild(option);
            });

            // Conditionally fetch courses if needed
            if (shouldFetchCourses) {
                console.log('Fetching courses because shouldFetchCourses is true');
                fetchCourses();
            }
        })
        .catch(error => {
            console.error('Error fetching semesters:', error);
            alert('An error occurred while fetching semesters');
        });
}

function fetchCourses() {
    const programId = document.getElementById('program_id').value;
    const semesterId = document.getElementById('semester_id').value;
    console.log('Fetching courses for program ID:', programId, 'and semester ID:', semesterId);

    fetch(`/getfiltercourses?program_id=${programId}&semester_id=${semesterId}`)
        .then(response => {
            console.log('Courses response status:', response.status);
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            const courseSelect = document.getElementById('course_id');
            courseSelect.innerHTML = ''; // Clear existing options

            // Add the default option
            const defaultOption = document.createElement('option');
            defaultOption.value = '';
            defaultOption.textContent = 'Select Course';
            courseSelect.appendChild(defaultOption);

            // Add new options
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

    // Initialize shouldFetchCourses based on the current page or context
    // Example: Set it to true if you need to fetch courses on this page
    shouldFetchCourses = false; // Set to false to prevent fetching courses on this page

    // Event listeners for dynamic updates
    document.getElementById('program_id').addEventListener('change', () => {
        fetchSemesters();
    });

    // Only add event listener for semester change if fetching courses is required
    if (shouldFetchCourses) {
        document.getElementById('semester_id').addEventListener('change', fetchCourses);
    }
});