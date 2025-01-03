// likes and  dislike buttons functionalities
document.addEventListener('DOMContentLoaded', () => {
    // Only add event listeners if user is authenticated
    const likeButtons = document.querySelectorAll('.btn-like');
    const dislikeButtons = document.querySelectorAll('.btn-dislike');

    if (likeButtons.length > 0) {
        // Handle Like Button Click
        likeButtons.forEach(button => {
            button.addEventListener('click', async () => {
                const postId = button.getAttribute('data-id');
                const promptContainer = document.getElementById(`reaction-prompt-${postId}`);
                
                // Check if user is authenticated by looking for user-info in the header
                const isAuthenticated = document.querySelector('.user-info') !== null;
                
                if (!isAuthenticated) {
                    // Show guest prompt if not authenticated
                    promptContainer.style.display = promptContainer.style.display === 'none' ? 'block' : 'none';
                    return;
                }

                try {
                    const response = await fetch(`/like-post?id=${postId}`, { method: 'POST' });
                    const data = await response.json();
                    
                    // Update all instances of this post's like/dislike counts
                    updateReactionCounts(postId, data.likes, data.dislikes);
                } catch (err) {
                    console.error('Error:', err);
                }
            });
        });

        // Handle Dislike Button Click
        dislikeButtons.forEach(button => {
            button.addEventListener('click', async () => {
                const postId = button.getAttribute('data-id');
                const promptContainer = document.getElementById(`reaction-prompt-${postId}`);
                
                // Check if user is authenticated by looking for user-info in the header
                const isAuthenticated = document.querySelector('.user-info') !== null;
                
                if (!isAuthenticated) {
                    // Show guest prompt if not authenticated
                    promptContainer.style.display = promptContainer.style.display === 'none' ? 'block' : 'none';
                    return;
                }

                try {
                    const response = await fetch(`/dislike-post?id=${postId}`, { method: 'POST' });
                    const data = await response.json();
                    
                    // Update all instances of this post's like/dislike counts
                    updateReactionCounts(postId, data.likes, data.dislikes);
                } catch (err) {
                    console.error('Error:', err);
                }
            });
        });

        // Close prompts when clicking elsewhere
        document.addEventListener('click', (event) => {
            if (!event.target.closest('.btn-like') && 
                !event.target.closest('.btn-dislike') && 
                !event.target.closest('.guest-reaction-prompt')) {
                document.querySelectorAll('.reaction-prompt-container').forEach(container => {
                    container.style.display = 'none';
                });
            }
        });
    }
});

// Function to update all instances of a post's reaction counts
function updateReactionCounts(postId, likes, dislikes) {
    // Update in dashboard view
    const dashboardLikeButtons = document.querySelectorAll(`.btn-like[data-id="${postId}"]`);
    const dashboardDislikeButtons = document.querySelectorAll(`.btn-dislike[data-id="${postId}"]`);

    dashboardLikeButtons.forEach(button => {
        button.querySelector('span').innerText = likes;
    });

    dashboardDislikeButtons.forEach(button => {
        button.querySelector('span').innerText = dislikes;
    });

    // Update in view-post view if we're on that page
    const viewPostLikeButton = document.querySelector(`.post-actions .btn-like[data-id="${postId}"] span`);
    const viewPostDislikeButton = document.querySelector(`.post-actions .btn-dislike[data-id="${postId}"] span`);

    if (viewPostLikeButton) viewPostLikeButton.innerText = likes;
    if (viewPostDislikeButton) viewPostDislikeButton.innerText = dislikes;
}
