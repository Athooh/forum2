document.addEventListener('DOMContentLoaded', () => {
    const commentButtons = document.querySelectorAll('.btn-comment');
    
    commentButtons.forEach(button => {
        button.addEventListener('click', () => {
            const postId = button.getAttribute('data-id');
            
            // Check if we're on the dashboard page
            const isDashboard = window.location.pathname === '/dashboard' || 
                              window.location.pathname === '/';
            
            if (isDashboard) {
                // Redirect to view-post page
                window.location.href = `/view-post?id=${postId}#comments-section`;
                return;
            }
            
            // Existing comment form logic for view-post page
            const commentForm = document.getElementById(`comment-form-${postId}`);
            const isAuthenticated = document.querySelector('.user-info') !== null;
            
            if (!isAuthenticated) {
                const guestPrompt = commentForm.querySelector('.guest-comment-prompt');
                if (guestPrompt) {
                    commentForm.style.display = commentForm.style.display === 'none' ? 'block' : 'none';
                }
                return;
            }
            
            if (commentForm.style.display === 'none') {
                document.querySelectorAll('.comment-form-container').forEach(form => {
                    form.style.display = 'none';
                });
                commentForm.style.display = 'block';
                const textarea = commentForm.querySelector('textarea');
                if (textarea) {
                    textarea.focus();
                }
            } else {
                commentForm.style.display = 'none';
            }
        });
    });

    // Handle form submission
    document.querySelectorAll('.quick-comment-form').forEach(form => {
        form.addEventListener('submit', async (e) => {
            e.preventDefault();
            
            // Check if user is authenticated
            const isAuthenticated = document.querySelector('.user-info') !== null;
            if (!isAuthenticated) {
                return;
            }
            
            const postId = form.querySelector('input[name="post_id"]').value;
            const comment = form.querySelector('textarea').value;
            
            try {
                const response = await fetch(`/add-comment?post_id=${postId}`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/x-www-form-urlencoded',
                    },
                    body: `comment=${encodeURIComponent(comment)}`
                });

                if (response.ok) {
                    // Update comment count
                    const countSpan = document.querySelector(`.btn-comment[data-id="${postId}"] span`);
                    const currentCount = parseInt(countSpan.textContent);
                    countSpan.textContent = currentCount + 1;
                    
                    // Clear and hide the form
                    form.reset();
                    form.parentElement.style.display = 'none';
                } else {
                    alert('Failed to post comment. Please try again.');
                }
            } catch (error) {
                console.error('Error posting comment:', error);
                alert('Failed to post comment. Please try again.');
            }
        });
    });
}); 