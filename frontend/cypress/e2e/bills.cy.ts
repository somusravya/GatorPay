describe('Bills Page', () => {
    beforeEach(() => {
        // Mock authentication state
        window.localStorage.setItem('gatorpay_token', 'fake-jwt-token');
        
        // Mock user profile response
        cy.intercept('GET', '**/api/v1/auth/me', {
            statusCode: 200,
            body: {
                success: true,
                data: {
                    user: { id: 'user-1', first_name: 'Test', last_name: 'User', email: 'test@example.com' },
                    wallet: { balance: 1540.50, currency: 'USD' }
                }
            }
        }).as('getMe');

        // Mock bill categories
        cy.intercept('GET', '**/api/v1/bills/categories', {
            statusCode: 200,
            body: {
                success: true,
                data: ["electricity", "internet", "phone", "water"]
            }
        }).as('getCategories');

        // Mock saved billers
        cy.intercept('GET', '**/api/v1/bills/saved', {
            statusCode: 200,
            body: {
                success: true,
                data: [
                    { id: '1', biller_id: 'biller-1', nickname: 'Home Internet', account_number: '123456789', biller: { name: 'Comcast', category: 'internet' } }
                ]
            }
        }).as('getSavedBillers');

        cy.visit('/bills');
    });

    it('should display the bill categories', () => {
        cy.wait(['@getMe', '@getCategories']);
        cy.contains(/electricity/i).should('be.visible');
        cy.contains(/internet/i).should('be.visible');
        cy.contains(/phone/i).should('be.visible');
    });

    it('should display saved billers', () => {
        cy.wait('@getSavedBillers');
        cy.contains(/Home Internet/i).should('be.visible');
        cy.contains(/Comcast/i).should('exist');
    });

    it('should allow navigation to pay a new bill', () => {
        // Just checking that elements exist to initiate a payment
        cy.get('button, a').contains(/pay|new/i).should('exist');
    });
});
