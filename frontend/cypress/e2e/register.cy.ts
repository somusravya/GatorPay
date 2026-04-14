describe('Register Page', () => {
    beforeEach(() => {
        cy.visit('/register');
    });

    it('should display the registration form with all required fields', () => {
        cy.get('input[formControlName="firstName"], input[name="firstName"], input[placeholder*="First Name" i]').should('be.visible');
        cy.get('input[formControlName="lastName"], input[name="lastName"], input[placeholder*="Last Name" i]').should('be.visible');
        cy.get('input[formControlName="username"], input[name="username"], input[placeholder*="Username" i]').should('be.visible');
        cy.get('input[type="email"], input[name="email"]').should('be.visible');
        cy.get('input[type="password"], input[name="password"]').should('be.visible');
        // Check for confirm password and phone based on standard forms
        cy.get('input[type="tel"], input[name="phone"], input[placeholder*="Phone" i]').should('exist');
    });

    it('should allow typing in the registration fields', () => {
        cy.get('input[formControlName="firstName"], input[name="firstName"], input[placeholder*="First Name" i]')
            .type('John');
        
        cy.get('input[formControlName="lastName"], input[name="lastName"], input[placeholder*="Last Name" i]')
            .type('Doe');

        cy.get('input[type="email"], input[name="email"]')
            .type('newuser@example.com');
            
        cy.get('input[type="password"], input[name="password"]').first()
            .type('TestPassword123!');
    });

    it('should show a link back to the login page', () => {
        cy.get('a[href*="login"], a[routerLink*="login"], button').contains(/log in|login|sign in/i).should('exist');
    });
});
