// password validation
document.getElementById("password").addEventListener("input", function () {
    console.log("Password input detected!");
    const password = this.value;
    const strengthMeter = document.getElementById("strength-meter");
    const strengthText = document.getElementById("strength-text");
    if (!strengthMeter || !strengthText) {
        console.error("Strength meter or text not found!");
    }

    let strength = 0;

    const regex = {
        length: /.{8,}/,
        upper: /[A-Z]/,
        lower: /[a-z]/,
        number: /\d/,
        special: /[!@#$%^&*(),.?":{}|<>]/,
    };

    // Check each criterion and increase strength accordingly
    if (regex.length.test(password)){
        console.log("Length criteria met");
        strength++;
    }
    if (regex.upper.test(password)){
        console.log("Uppercase criteria met");
        strength++;
    } 
    if (regex.lower.test(password)){
        console.log("Lower criteria met");
        strength++;
    } 
    if (regex.number.test(password)){
        console.log("Number criteria met");
        strength++;
    } 
    if (regex.special.test(password)){
        console.log("Special criteria met");
        strength++;
    } 

    // Update strength meter
    const meterWidth = (strength / 5) * 100;
    strengthMeter.style.width = meterWidth + "%";

    // Update text and color based on strength
    if (strength === 1) {
        strengthText.textContent = "Weak";
        strengthMeter.style.backgroundColor = "red";
    } else if (strength === 2) {
        strengthText.textContent = "Fair";
        strengthMeter.style.backgroundColor = "orange";
    } else if (strength === 3) {
        strengthText.textContent = "Good";
        strengthMeter.style.backgroundColor = "yellow";
    } else if (strength === 4) {
        strengthText.textContent = "Strong";
        strengthMeter.style.backgroundColor = "green";
    } else if (strength === 5) {
        strengthText.textContent = "Very Strong";
        strengthMeter.style.backgroundColor = "darkgreen";
    }
});
