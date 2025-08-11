import { Component, ElementRef, OnInit, ViewChild } from '@angular/core';
import { FormsModule } from '@angular/forms'; // Import FormsModule
import { CommonModule } from '@angular/common'; // Import CommonModule
import { ApiService } from '../../services/api.service';
import { AuthenticationService } from '../../services/authentication.service';
import { ClarityModule } from '@clr/angular';
import { ConfigService } from 'src/app/services/config.service';

@Component({
  selector: 'app-billing',
  standalone: true,
  imports: [FormsModule, CommonModule, ClarityModule], // Add CommonModule here
  templateUrl: './billing.component.html',
  styleUrls: ['./billing.component.scss']
})
export class BillingComponent implements OnInit {
  @ViewChild('generateChartsLoading') generateChartsLoading!: ElementRef;
  @ViewChild('generateChartsButton') generateChartsButton!: ElementRef;
  @ViewChild('overcommitLoading') overcommitLoading!: ElementRef;
  @ViewChild('overcommitButton') overcommitButton!: ElementRef;

  chartYear = '';
  chartMonth = '';
  overcommitYear = '';
  overcommitMonth = '';
  selectedCustomerID = '';
  billingApiEndpoint: string = ''; // Initialize with an empty string
  limit = '';
  customers: { id: string[]; name: string; overcommitment?: string }[] = [];
  routes: string[] = [];
  tasks: string[] = []; // Add a property to store tasks

  constructor(
    private apiService: ApiService,
    private authService: AuthenticationService,
    private configService: ConfigService // Inject AuthenticationService
  ) {
    this.billingApiEndpoint = this.configService.getEnv('billingApiEndpoint'); // Read billingApiEndpoint from ConfigService
  }

  ngOnInit(): void {
    this.fetchCustomers();
    this.fetchRoutes();
  }


  fetchCustomers(): void {
    this.apiService.getCustomers().subscribe({
      next: (data: any) => {
        // Transform the API response into the `customers` array
        this.customers = Object.entries(data).map(([name, ids]) => ({
          id: ids as string[],
          name
        }));
      },
      error: (error) => {
        console.error('Error fetching customers:', error);
      }
    });
  }

  fetchRoutes(): void {
    this.apiService.getRoutes().subscribe({
      next: (data: string[]) => {
        // Assign the API response directly to the `routes` array
        this.routes = data;
      },
      error: (error) => {
        console.error('Error fetching routes:', error);
      }
    });
  }

  onSubmitGenerateCharts(event: Event): void {
    event.preventDefault();
    this.submitGenerateCharts(this.chartYear, this.chartMonth);
  }

  onSubmitOvercommitment(event: Event): void {
    event.preventDefault();
    this.submitOvercommitment(this.overcommitYear, this.overcommitMonth, this.selectedCustomerID, this.limit);
  }

  submitGenerateCharts(year: string, month: string): void {
    const loading = this.generateChartsLoading.nativeElement;
    const button = this.generateChartsButton.nativeElement;

    loading.style.display = 'block';
    button.disabled = true;

    this.apiService.generateCharts({ year, month }).subscribe({
      next: (response: any) => {
        loading.style.display = 'none';
        button.disabled = false;
        location.reload();
      },
      error: (error) => {
        loading.style.display = 'none';
        button.disabled = false;
        alert('Error generating charts: ' + error.message);
      }
    });
  }

  submitOvercommitment(year: string, month: string, customerID: string, limit: string): void {
    const loading = this.overcommitLoading.nativeElement;
    const button = this.overcommitButton.nativeElement;

    loading.style.display = 'block';
    button.disabled = true;

    this.apiService.calculateOvercommitment({ year, month, customerID, limit }).subscribe({
      next: (data: any) => {
        loading.style.display = 'none';
        button.disabled = false;

        if (data.error) {
          alert('Error calculating overcommitment: ' + data.error);
        } else {
          console.log('Customer ID:', customerID,'Overcommitment data:', data.overcommitment);


          const customer = this.customers.find(c => c.id.includes(customerID.split(',')[0]));
          if (customer) {
            customer.overcommitment = data.overcommitment === 0 ? '0' : data.overcommitment; // Ensure "0" is written as a string
            // alert(`Overcommitment for ${customer.name} in ${year}-${month}: ${customer.overcommitment}`);
          }
        }
      },
      error: (error) => {
        loading.style.display = 'none';
        button.disabled = false;
        alert('Error: ' + error.message);
      }
    });
  }

  openChart(route: string): void {
    const authToken = this.authService['oauthService'].getAccessToken(); // Get the token from AuthenticationService
    const url = this.billingApiEndpoint + route;

    // Open the chart in a new tab with the Authorization header
    const headers = new Headers();
    headers.append('Authorization', `Bearer ${authToken}`);

    fetch(url, { headers })
      .then((response) => {
        if (response.ok) {
          return response.blob();
        } else {
          throw new Error('Failed to fetch chart');
        }
      })
      .then((blob) => {
        const blobUrl = URL.createObjectURL(blob);
        window.open(blobUrl, '_blank'); // Open the chart in a new tab
      })
      .catch((error) => {
        console.error('Error opening chart:', error);
      });
  }
}
