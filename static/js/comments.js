document.addEventListener('DOMContentLoaded', () => {
    const commentButtons = document.querySelectorAll('.btn-comment');
    
    commentButtons.forEach(button => {
        button.addEventListener('click', () => {
            const postId = button.getAttribute('data-id');
            const commentForm = document.getElementById(`comment-form-${postId}`);
            
            // Toggle the comment form visibility
            if (commentForm.style.display === 'none') {
                // Hide all other open comment forms first
                document.querySelectorAll('.comment-form-container').forEach(form => {
                    form.style.display = 'none';
                });
                // Show this comment form
                commentForm.style.display = 'block';
                // Focus on the textarea
                commentForm.querySelector('textarea').focus();
            } else {
                commentForm.style.display = 'none';
            }
        });
    });

    // Handle form submission
    document.querySelectorAll('.quick-comment-form').forEach(form => {
        form.addEventListener('submit', async (e) => {
            e.preventDefault();
            
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